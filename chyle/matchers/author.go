package matchers

import (
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"
)

// author is commit author matcher
type author struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit author field
func (a author) Match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

func buildAuthor(re *regexp.Regexp) Matcher {
	return author{re}
}
