package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jymannob/namecheck"
	_ "github.com/jymannob/namecheck/facebook"
	_ "github.com/jymannob/namecheck/github"
	_ "github.com/jymannob/namecheck/instagram"
)

type result struct {
	username  string
	platform  string
	valid     bool
	available bool
	err       error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: skdhsdkfjhsdf")
		os.Exit(1)
	}
	username := os.Args[1]
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
		fmt.Printf("%+v\n", r)
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
