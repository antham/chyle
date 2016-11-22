package chyle

import (
	"regexp"

	"gopkg.in/src-d/go-git.v4"
)

// Matcher describe a way of applying a matcher against a commit
type Matcher interface {
	Match(*git.Commit) bool
}

// MergeCommitMatcher match merge commit message
type MergeCommitMatcher struct {
}

// Match is valid if commit is a merge commit
func (m MergeCommitMatcher) Match(commit *git.Commit) bool {
	return commit.NumParents() == 2
}

// RegularCommitMatcher match regular commit message
type RegularCommitMatcher struct {
}

// Match is valid if commit is not a merge commit
func (r RegularCommitMatcher) Match(commit *git.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

// MessageMatcher is commit message matcher
type MessageMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit message
func (m MessageMatcher) Match(commit *git.Commit) bool {
	return m.regexp.MatchString(commit.Message)
}

// CommitterMatcher is commit committer matcher
type CommitterMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit committer field
func (c CommitterMatcher) Match(commit *git.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

// AuthorMatcher is commit author matcher
type AuthorMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit author field
func (a AuthorMatcher) Match(commit *git.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
}

// Filter commits that don't fit any matchers
func Filter(matchers *[]Matcher, commits *[]git.Commit) *[]git.Commit {
	results := []git.Commit{}

	for _, commit := range *commits {
		add := true
		for _, matcher := range *matchers {
			if !matcher.Match(&commit) {
				add = false
			}
		}

		if add {
			results = append(results, commit)
		}
	}

	return &results
}

// TransformCommitsToMap extract useful commits data in hash map table
func TransformCommitsToMap(commits *[]git.Commit) *[]map[string]interface{} {
	commitMaps := []map[string]interface{}{}

	for _, c := range *commits {
		commitMap := map[string]interface{}{
			"id":             c.ID().String(),
			"authorName":     c.Author.Name,
			"authorEmail":    c.Author.Email,
			"authorDate":     c.Author.When.String(),
			"committerName":  c.Committer.Name,
			"committerEmail": c.Committer.Email,
			"committerDate":  c.Committer.When.String(),
			"message":        c.Message,
			"isMerge":        c.NumParents() == 2,
		}

		commitMaps = append(commitMaps, commitMap)
	}

	return &commitMaps
}
