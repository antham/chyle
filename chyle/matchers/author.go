package matchers

import (
	"regexp"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// author is commit author matcher
type author struct {
	regexp *regexp.Regexp
}

func (a author) Match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

func newAuthor(re *regexp.Regexp) Matcher {
	return author{re}
}
