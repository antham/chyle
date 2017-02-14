package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// AuthorMatcher is commit author matcher
type AuthorMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit author field
func (a AuthorMatcher) Match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

func buildAuthorMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return AuthorMatcher{r}, nil
}
