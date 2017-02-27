package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// MessageMatcher is commit message matcher
type MessageMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit message
func (m MessageMatcher) Match(commit *object.Commit) bool {
	return m.regexp.MatchString(commit.Message)
}

func buildMessageMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return MessageMatcher{r}, nil
}

// removePGPKey fix library issue that don't trim PGP key from message
func removePGPKey(message string) string {
	return regexp.MustCompile("(?s).*?-----END PGP SIGNATURE-----\n\n").ReplaceAllString(message, "")
}
