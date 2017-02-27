package chyle

import (
	"fmt"

	"srcd.works/go-git.v4/plumbing/object"
)

// mergeCommitMatcher match merge commit message
type mergeCommitMatcher struct {
}

// match is valid if commit is a merge commit
func (m mergeCommitMatcher) match(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// regularCommitMatcher match regular commit message
type regularCommitMatcher struct {
}

// match is valid if commit is not a merge commit
func (r regularCommitMatcher) match(commit *object.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

func buildTypeMatcher(key string, value string) (matcher, error) {
	switch value {
	case "regular":
		return regularCommitMatcher{}, nil
	case "merge":
		return mergeCommitMatcher{}, nil
	}

	return nil, fmt.Errorf(`"%s" must be "regular" or "merge", "%s" given`, key, value)
}

func solveType(commit *object.Commit) string {
	switch commit.NumParents() {
	case 0, 1:
		return "regular"
	case 2:
		return "merge"
	default:
		return "unknown"
	}
}
