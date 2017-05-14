package matchers

import (
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// message is commit message matcher
type message struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit message
func (m message) Match(commit *object.Commit) bool {
	return m.regexp.MatchString(commit.Message)
}

func buildMessage(re *regexp.Regexp) Matcher {
	return message{re}
}

// removePGPKey fix library issue that don't trim PGP key from message
func removePGPKey(message string) string {
	return regexp.MustCompile("(?s).*?-----END PGP SIGNATURE-----\n\n").ReplaceAllString(message, "")
}
