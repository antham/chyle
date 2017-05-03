package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTypeMatcher(t *testing.T) {
	assert.Equal(t, regularCommitMatcher{}, buildTypeMatcher(regularTypeMatcher))
	assert.Equal(t, mergeCommitMatcher{}, buildTypeMatcher(mergeTypeMatcher))
}

func TestSolveType(t *testing.T) {
	c := getCommitFromRef("HEAD")

	assert.Equal(t, 1, c.NumParents())
	assert.Equal(t, regularTypeMatcher, solveType(c))

	c = getCommitFromRef("HEAD~2")

	assert.Equal(t, 2, c.NumParents())
	assert.Equal(t, mergeTypeMatcher, solveType(c))

	c = getCommitFromRef("HEAD~4")

	assert.Equal(t, 0, c.NumParents())
	assert.Equal(t, regularTypeMatcher, solveType(c))
}
