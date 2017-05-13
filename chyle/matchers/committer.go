package matchers

import (
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// committerMatcher is commit committer matcher
type committerMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit committer field
func (c committerMatcher) Match(commit *object.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

func buildCommitterMatcher(value string) Matcher {
	return committerMatcher{regexp.MustCompile(value)}
}
