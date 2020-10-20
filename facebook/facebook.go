package facebook

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jymannob/namecheck"
)

const (
	platform         = "Facebook"
	minLen           = 5
	maxLen           = 255
	illegalSuffix    = ".net"
	illegalSubstring = ".com"
	endpointTmpl     = "https://www.facebook.com/%s"
)

var legalRegexp = regexp.MustCompile("^[0-9A-Za-z.]*$")

type Facebook struct {
	Client namecheck.Client
}

func init() {
	gh := Facebook{
		Client: http.DefaultClient,
	}
	namecheck.Register(&gh)
}

func (g *Facebook) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		onlyContainsLegalChars(username) &&
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

func containsNoIllegalSuffix(username string) bool {
	return !strings.HasSuffix(username, illegalSuffix)
}

func (g *Facebook) IsAvailable(username string) (bool, error) {
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

func (_ *Facebook) String() string {
	return platform
}
