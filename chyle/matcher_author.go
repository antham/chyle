package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// authorMatcher is commit author matcher
type authorMatcher struct {
	regexp *regexp.Regexp
}

// match apply a regexp against commit author field
func (a authorMatcher) match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

func buildauthorMatcher(key string, value string) (matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return authorMatcher{r}, nil
}
