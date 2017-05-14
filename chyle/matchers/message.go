package matchers

import (
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// messageMatcher is commit message matcher
type messageMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit message
func (m messageMatcher) Match(commit *object.Commit) bool {
	return m.regexp.MatchString(commit.Message)
}

func buildMessageMatcher(re *regexp.Regexp) Matcher {
	return messageMatcher{re}
}

// removePGPKey fix library issue that don't trim PGP key from message
func removePGPKey(message string) string {
	return regexp.MustCompile("(?s).*?-----END PGP SIGNATURE-----\n\n").ReplaceAllString(message, "")
}
