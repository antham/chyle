package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvDecorator(t *testing.T) {
	setenv("TESTENVDECORATOR", "this is a test")

	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASENVDECORATOR = true
	chyleConfig.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"WHATEVER": {
			"envDecoratorTesting",
			"TESTENVDECORATOR",
		},
	}

	metadatas := map[string]interface{}{}

	e := buildEnvDecorators()
	m, err := e[0].decorate(&metadatas)

	assert.NoError(t, err, "Must returns no errors")
	assert.Equal(t, map[string]interface{}{"envDecoratorTesting": "this is a test"}, *m, "Must dump environment variable in given destination key")
}
