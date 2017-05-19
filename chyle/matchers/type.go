package matchers

import (
	"srcd.works/go-git.v4/plumbing/object"
)

const (
	regularType = "regular"
	mergeType   = "merge"
)

// mergeCommit match merge commit message
type mergeCommit struct {
}

// match is valid if commit is a merge commit
func (m mergeCommit) Match(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// regularCommit match regular commit message
type regularCommit struct {
}

// Match is valid if commit is not a merge commit
func (r regularCommit) Match(commit *object.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

func newType(key string) Matcher {
	if key == regularType {
		return regularCommit{}
	}

	return mergeCommit{}
}

func solveType(commit *object.Commit) string {
	if commit.NumParents() == 2 {
		return mergeType
	}

	return regularType
}

// GetTypes returns all defined matchers types
func GetTypes() []string {
	return []string{regularType, mergeType}
}
