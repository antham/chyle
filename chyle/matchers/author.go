package matchers

import (
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// authorMatcher is commit author matcher
type authorMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit author field
func (a authorMatcher) Match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

func buildAuthorMatcher(re *regexp.Regexp) Matcher {
	return authorMatcher{re}
}
