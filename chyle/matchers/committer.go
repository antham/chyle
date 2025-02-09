package matchers

import (
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// committer is commit committer matcher
type committer struct {
	regexp *regexp.Regexp
}

func (c committer) Match(commit *object.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

func newCommitter(re *regexp.Regexp) Matcher {
	return committer{re}
}
