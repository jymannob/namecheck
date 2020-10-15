package namecheck

import "fmt"

var checkers []Checker

type Validator interface {
	IsValid(username string) bool
}

type Availabler interface {
	IsAvailable(username string) (bool, error)
}

type Checker interface {
	fmt.Stringer
	Validator
	Availabler
}

func Register(c Checker) {
	checkers = append(checkers, c)
}

func Checkers() []Checker {
	return checkers
}
