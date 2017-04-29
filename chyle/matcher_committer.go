package chyle

import (
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

func buildCommitterMatcher(value string) matcher {
	return committerMatcher{regexp.MustCompile(value)}
}
