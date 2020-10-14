package github

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	minLen           = 1
	maxLen           = 39
	illegalPrefix    = "-"
	illegalSuffix    = "-"
	illegalSubstring = "--"
)

var legalRegexp = regexp.MustCompile("^[-0-9A-Za-z]*$")

func IsValid(username string) bool {
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
