package facebook

import (
	"strings"
	"testing"

	"github.com/jub0bs/namecheck"
)

var (
	_  namecheck.Checker = (*Facebook)(nil)
	fb *Facebook
)

func TestValidateFailsOnNamesThatContainIllegalChars(t *testing.T) {
	username := "underscore_"
	want := false
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalSuffix(t *testing.T) {
	username := "notok-"
	want := false
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatContainIllegalSubstring(t *testing.T) {
	username := "no.comok"
	want := false
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooShort(t *testing.T) {
	username := ""
	want := false
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreLongEnough(t *testing.T) {
	username := "aaaaa"
	want := true
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateFailsOnNamesThatAreTooLong(t *testing.T) {
	username := strings.Repeat("a", maxLen+1)
	want := false
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}

func TestValidateSucceedsOnNamesThatAreShortEnough(t *testing.T) {
	username := strings.Repeat("a", maxLen)
	want := true
	got := fb.IsValid(username)
	if got != want {
		t.Errorf("IsValid(%s) = %t; want %t", username, got, want)
	}
}
