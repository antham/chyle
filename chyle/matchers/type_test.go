package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildType(t *testing.T) {
	assert.Equal(t, regularCommit{}, buildType(regularType))
	assert.Equal(t, mergeCommit{}, buildType(mergeType))
}
