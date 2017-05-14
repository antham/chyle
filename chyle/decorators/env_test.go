package decorators

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvDecorator(t *testing.T) {
	err := os.Setenv("TESTENVDECORATOR", "this is a test")

	assert.NoError(t, err)

	envs := map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"WHATEVER": {
			"envDecoratorTesting",
			"TESTENVDECORATOR",
		},
	}

	metadatas := map[string]interface{}{}

	e := buildEnvDecorators(envs)
	m, err := e[0].Decorate(&metadatas)

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"envDecoratorTesting": "this is a test"}, *m, "Must dump environment variable in given destination key")
}
