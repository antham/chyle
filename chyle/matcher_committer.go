package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// committerMatcher is commit committer matcher
type committerMatcher struct {
	regexp *regexp.Regexp
}

// match apply a regexp against commit committer field
func (c committerMatcher) match(commit *object.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

func buildCommitterMatcher(key string, value string) (matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return committerMatcher{r}, nil
}
