package github

import (
	"strings"
	"testing"

	"github.com/jub0bs/namecheck"
)

var (
	_  namecheck.Checker = (*GitHub)(nil)
	gh *GitHub
)

func TestValidateFailsOnNamesThatContainIllegalChars(t *testing.T) {
	username := "underscore_"
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalPrefix(t *testing.T) {
	username := "-notok"
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalSuffix(t *testing.T) {
	username := "notok-"
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalSubstring(t *testing.T) {
	username := "no--ok"
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooShort(t *testing.T) {
	username := ""
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreLongEnough(t *testing.T) {
	username := "a"
	want := true
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooLong(t *testing.T) {
	username := strings.Repeat("a", maxLen+1)
	want := false
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreShortEnough(t *testing.T) {
	username := strings.Repeat("a", maxLen)
	want := true
	got := gh.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}
