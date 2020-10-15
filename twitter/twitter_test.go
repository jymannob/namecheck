package twitter

import (
	"strings"
	"testing"

	"github.com/jub0bs/namecheck"
)

var (
	_  namecheck.Checker = (*Twitter)(nil)
	tw *Twitter
)

func TestValidateFailsOnNamesThatContainIllegalChars(t *testing.T) {
	username := "hyphen-"
	want := false
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalPattern(t *testing.T) {
	username := "fooTwItterbar"
	want := false
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooShort(t *testing.T) {
	username := ""
	want := false
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreLongEnough(t *testing.T) {
	username := "a"
	want := true
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooLong(t *testing.T) {
	username := strings.Repeat("a", maxLen+1)
	want := false
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreShortEnough(t *testing.T) {
	username := strings.Repeat("a", maxLen)
	want := true
	got := tw.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}
