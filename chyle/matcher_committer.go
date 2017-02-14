package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// CommitterMatcher is commit committer matcher
type CommitterMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit committer field
func (c CommitterMatcher) Match(commit *object.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

func buildCommitterMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return CommitterMatcher{r}, nil
}
