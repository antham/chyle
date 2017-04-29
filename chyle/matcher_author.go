package chyle

import (
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

func buildAuthorMatcher(value string) matcher {
	return authorMatcher{regexp.MustCompile(value)}
}
