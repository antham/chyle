package matchers

import (
	"srcd.works/go-git.v4/plumbing/object"
)

const (
	regularTypeMatcher = "regular"
	mergeTypeMatcher   = "merge"
)

// mergeCommitMatcher match merge commit message
type mergeCommitMatcher struct {
}

// match is valid if commit is a merge commit
func (m mergeCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// regularCommitMatcher match regular commit message
type regularCommitMatcher struct {
}

// Match is valid if commit is not a merge commit
func (r regularCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

func buildTypeMatcher(key string) Matcher {
	if key == regularTypeMatcher {
		return regularCommitMatcher{}
	}

	return mergeCommitMatcher{}
}

func solveType(commit *object.Commit) string {
	if commit.NumParents() == 2 {
		return mergeTypeMatcher
	}

	return regularTypeMatcher
}

// GetTypes returns all defined matchers types
func GetTypes() []string {
	return []string{regularTypeMatcher, mergeTypeMatcher}
}
