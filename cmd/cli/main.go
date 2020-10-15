package main

import (
	"fmt"
	"os"

	"github.com/jub0bs/namecheck"
	_ "github.com/jub0bs/namecheck/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: skdhsdkfjhsdf")
		os.Exit(1)
	}
	username := os.Args[1]
	for _, c := range namecheck.Checkers() {
		valid := c.IsValid(username)
		if !valid {
			fmt.Printf("%q is invalid on %v\n", username, c)
			continue
		}
		avail, err := c.IsAvailable(username)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !avail {
			fmt.Printf("%q is unavailable on %v\n", username, c)
			continue
		}
		fmt.Printf("%q is available on %v\n", username, c)
	}
}
