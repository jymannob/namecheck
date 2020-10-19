package github

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/mockclient"
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

func TestIsAvailable(t *testing.T) {
	cases := []struct {
		label         string
		username      string
		client        namecheck.Client
		available     bool
		errorOccurred bool
	}{
		{
			label:     "notfound",
			username:  "dummy",
			client:    mockclient.WithStatusCode(http.StatusNotFound),
			available: true,
		}, {
			label:    "ok",
			username: "dummy",
			client:   mockclient.WithStatusCode(http.StatusOK),
		}, {
			label:    "other", // other than 200, 404
			username: "dummy",
			client:   mockclient.WithStatusCode(999),
		}, {
			label:         "clienterror",
			username:      "dummy",
			client:        mockclient.WithError(errors.New("some network error")),
			available:     false,
			errorOccurred: true,
		},
	}

	const template = "IsAvailable(%q): got %t (and %s error); want %t (and %s error)"
	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			gh := GitHub{
				Client: c.client,
			}
			available, err := gh.IsAvailable(c.username)
			if available != c.available || (err != nil != c.errorOccurred) {
				t.Errorf(
					template,
					c.username,
					available,
					errorMsgHelper(err != nil),
					c.available,
					errorMsgHelper(c.errorOccurred))
			}
		})
	}
}

func errorMsgHelper(errorOccurred bool) string {
	if errorOccurred {
		return "some"
	}
	return "no"
}
