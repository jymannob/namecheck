package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/jub0bs/namecheck"
	_ "github.com/jub0bs/namecheck/github"
)

var count uint64
var m = make(map[string]uint64)
var mu sync.Mutex

type result struct {
	username  string
	platform  string
	valid     bool
	available bool
	err       error
}

func main() {
	http.Handle("/check", http.HandlerFunc(handle))
	http.Handle("/count", http.HandlerFunc(handleCount))
	http.Handle("/details", http.HandlerFunc(handleDetails))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	msg := fmt.Sprintf("%d", atomic.LoadUint64(&count))
	fmt.Fprintf(w, msg)
}

func handleDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	var msg string
	mu.Lock()
	{
		msg = fmt.Sprintf("%v", m)
	}
	mu.Unlock()
	fmt.Fprintf(w, msg)
}

func handle(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&count, 1)
	w.Header().Add("Content-Type", "text/plain")
	username := r.URL.Query().Get("username")
	mu.Lock()
	{
		m[username]++
	}
	mu.Unlock()
	if len(username) == 0 {
		http.Error(w, "missing 'username' query parameter", http.StatusBadRequest)
		return
	}
	ch := make(chan result)
	var wg sync.WaitGroup
	checkers := namecheck.Checkers()
	wg.Add(len(checkers))
	for _, c := range checkers {
		go check(c, username, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for res := range ch {
		fmt.Fprintf(w, "%+v\n", res)
	}
}

func check(
	c namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	ch chan<- result) {
	defer wg.Done()
	res := result{
		username: username,
		platform: c.String(),
	}
	valid := c.IsValid(username)
	if !valid {
		ch <- res
		return
	}
	res.valid = true
	avail, err := c.IsAvailable(username)
	if err != nil {
		res.err = err
		ch <- res
		return
	}
	if !avail {
		ch <- res
		return
	}
	res.available = true
	ch <- res
}
