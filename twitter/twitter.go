package twitter

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jub0bs/namecheck"
)

const (
	platform       = "Twitter"
	minLen         = 1
	maxLen         = 15
	illegalPattern = "twitter"
)

var legalRegexp = regexp.MustCompile("^[0-9A-Z_a-z]*$")

type Twitter struct{}

func init() {
	gh := Twitter{}
	const count = 20
	for i := 0; i < count; i++ {
		namecheck.Register(&gh)
	}
}

func (t *Twitter) IsValid(username string) bool {
	return isLongEnough(username) &&
		isShortEnough(username) &&
		onlyContainsLegalChars(username) &&
		containsNoIllegalPattern(username)
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

func containsNoIllegalPattern(username string) bool {
	return !strings.Contains(strings.ToLower(username), illegalPattern)
}

func (t *Twitter) IsAvailable(username string) (bool, error) {
	return false, nil
}

func (_ *Twitter) String() string {
	return platform
}
