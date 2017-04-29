package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvDecorator(t *testing.T) {
	setenv("DECORATORS_ENV_WHATEVER_VALUE", "hello world !")
	setenv("DECORATORS_ENV_WHATEVER_DESTKEY", "helloWorld")

	metadatas := map[string]interface{}{}

	e := envDecorator{"DECORATORS_ENV_WHATEVER_VALUE", "helloWorld"}
	m, err := e.decorate(&metadatas)

	assert.NoError(t, err, "Must returns no errors")
	assert.Equal(t, map[string]interface{}{"helloWorld": "hello world !"}, *m, "Must dump environment variable in given destination key")
}
