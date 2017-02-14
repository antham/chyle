package chyle

import (
	"fmt"

	"srcd.works/go-git.v4/plumbing/object"
)

// MergeCommitMatcher match merge commit message
type MergeCommitMatcher struct {
}

// Match is valid if commit is a merge commit
func (m MergeCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// RegularCommitMatcher match regular commit message
type RegularCommitMatcher struct {
}

// Match is valid if commit is not a merge commit
func (r RegularCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

func buildTypeMatcher(key string, value string) (Matcher, error) {
	switch value {
	case "regular":
		return RegularCommitMatcher{}, nil
	case "merge":
		return MergeCommitMatcher{}, nil
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
