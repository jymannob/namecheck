package instagram

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode/utf8"

	"github.com/jymannob/namecheck"
)

const (
	platform     = "Instagram"
	minLen       = 1
	maxLen       = 30
	endpointTmpl = "https://www.instagram.com/%s"
)

var legalRegexp = regexp.MustCompile("^[0-9A-Za-z._]*$")

type Instagram struct {
	Client namecheck.Client
}

func init() {
	gh := Instagram{
		Client: http.DefaultClient,
	}
	namecheck.Register(&gh)
}

func (g *Instagram) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		onlyContainsLegalChars(username)
}

func isLongEnough(username string) bool {
	return minLen <= utf8.RuneCountInString(username)
}

func isShortEnough(username string) bool {
	return utf8.RuneCountInString(username) <= maxLen
}

func onlyContainsLegalChars(username string) bool {
	return legalRegexp.MatchString(username)
}

func (g *Instagram) IsAvailable(username string) (bool, error) {
	endpoint := fmt.Sprintf(endpointTmpl, url.PathEscape(username))
	resp, err := g.Client.Get(endpoint)
	if err != nil {
		hlErr := namecheck.ErrUnknownAvailability{
			Username: username,
			Platform: platform,
			Cause:    err,
		}
		return false, &hlErr
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusNotFound, nil
}

func (_ *Instagram) String() string {
	return platform
}
