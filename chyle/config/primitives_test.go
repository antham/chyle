package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingEnvError(t *testing.T) {
	errorWithOneVarMissing := MissingEnvError{
		[]string{"TEST_1"},
	}

	assert.Equal(t, errorWithOneVarMissing.Envs(), []string{"TEST_1"})
	assert.EqualError(t, errorWithOneVarMissing, `environment variable missing : "TEST_1"`)

	errorWithSeveralVarsMissing := MissingEnvError{
		[]string{"TEST_1", "TEST_2"},
	}

	assert.Equal(t, errorWithSeveralVarsMissing.Envs(), []string{"TEST_1", "TEST_2"})
	assert.EqualError(t, errorWithSeveralVarsMissing, `environments variables missing : "TEST_1", "TEST_2"`)
}
