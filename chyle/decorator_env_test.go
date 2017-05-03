package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvDecorator(t *testing.T) {
	setenv("DECORATORS_ENV_WHATEVER_VALUE", "hello world !")
	setenv("DECORATORS_ENV_WHATEVER_DESTKEY", "helloWorld")

	chyleConfig.DECORATORS.ENV = map[string]map[string]string{"WHATEVER": {"VALUE": "DECORATORS_ENV_WHATEVER_VALUE", "DESTKEY": "helloWorld"}}

	metadatas := map[string]interface{}{}

	e := buildEnvDecorators()
	m, err := e[0].decorate(&metadatas)

	assert.NoError(t, err, "Must returns no errors")
	assert.Equal(t, map[string]interface{}{"helloWorld": "hello world !"}, *m, "Must dump environment variable in given destination key")
}
