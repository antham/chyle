package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	_, err := buildProcess()

	assert.NoError(t, err, "Must produces no errors")
}
