package main

import (
	"fmt"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/twitter"
)

func main() {
	const username = "jub0bs"
	if twitter.IsValid(username) && github.IsValid(username) {
		const fmtstr = "%s is valid on both Twitter and GitHub!\n"
		fmt.Printf(fmtstr, username)
	}
}
