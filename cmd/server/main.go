package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jub0bs/namecheck"
	_ "github.com/jub0bs/namecheck/github"
)

type result struct {
	username  string
	platform  string
	valid     bool
	available bool
	err       error
}

func main() {
	http.Handle("/", http.HandlerFunc(handle))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	username := r.URL.Query().Get("username")
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
	for r := range ch {
		fmt.Fprintf(w, "%+v\n", r)
	}
}

func check(
	c namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	ch chan<- result) {
	defer wg.Done()
	r := result{
		username: username,
		platform: c.String(),
	}
	valid := c.IsValid(username)
	if !valid {
		ch <- r
		return
	}
	r.valid = true
	avail, err := c.IsAvailable(username)
	if err != nil {
		r.err = err
		ch <- r
		return
	}
	if !avail {
		ch <- r
		return
	}
	r.available = true
	ch <- r
}
