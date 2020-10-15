package github

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jub0bs/namecheck"
)

const (
	platform         = "GitHub"
	minLen           = 1
	maxLen           = 39
	illegalPrefix    = "-"
	illegalSuffix    = "-"
	illegalSubstring = "--"
	endpointTmpl     = "https://github.com/%s"
)

var legalRegexp = regexp.MustCompile("^[-0-9A-Za-z]*$")

type GitHub struct {
	Client namecheck.Client
}

func init() {
	gh := GitHub{
		Client: http.DefaultClient,
	}
	const count = 20
	for i := 0; i < count; i++ {
		namecheck.Register(&gh)
	}
}

func (g *GitHub) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		onlyContainsLegalChars(username) &&
		containsNoIllegalPrefix(username) &&
		containsNoIllegalSuffix(username) &&
		containsNoIllegalSubstring(username)
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

func containsNoIllegalSubstring(username string) bool {
	return !strings.Contains(username, illegalSubstring)
}

func containsNoIllegalPrefix(username string) bool {
	return !strings.HasPrefix(username, illegalPrefix)
}

func containsNoIllegalSuffix(username string) bool {
	return !strings.HasSuffix(username, illegalSuffix)
}

func (g *GitHub) IsAvailable(username string) (bool, error) {
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

func (_ *GitHub) String() string {
	return platform
}
