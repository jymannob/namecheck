package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jub0bs/namecheck"
	_ "github.com/jub0bs/namecheck/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: skdhsdkfjhsdf")
		os.Exit(1)
	}
	username := os.Args[1]
	var wg sync.WaitGroup
	checkers := namecheck.Checkers()
	wg.Add(len(checkers))
	for _, c := range namecheck.Checkers() {
		go check(c, username, &wg)
	}
	wg.Wait()
}

func check(c namecheck.Checker, username string, wg *sync.WaitGroup) {
	defer wg.Done()
	valid := c.IsValid(username)
	if !valid {
		fmt.Printf("%q is invalid on %v\n", username, c)
		return
	}
	avail, err := c.IsAvailable(username)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !avail {
		fmt.Printf("%q is unavailable on %v\n", username, c)
		return
	}
	fmt.Printf("%q is available on %v\n", username, c)
}
