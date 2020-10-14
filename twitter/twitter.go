package twitter

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	minLen         = 1
	maxLen         = 15
	illegalPattern = "twitter"
)

var legalRegexp = regexp.MustCompile("^[0-9A-Z_a-z]*$")

func IsValid(username string) bool {
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
