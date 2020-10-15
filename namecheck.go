package namecheck

import "fmt"

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
