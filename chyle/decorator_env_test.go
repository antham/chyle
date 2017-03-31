package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestEnvDecorator(t *testing.T) {
	setenv("DECORATORS_ENV_WHATEVER_VALUE", "hello world !")
	setenv("DECORATORS_ENV_WHATEVER_DESTKEY", "helloWorld")

	metadatas := map[string]interface{}{}

	e := envDecorator{"DECORATORS_ENV_WHATEVER_VALUE", "DECORATORS_ENV_WHATEVER_DESTKEY"}
	m, err := e.decorate(&metadatas)

	assert.NoError(t, err, "Must returns no errors")
	assert.Equal(t, map[string]interface{}{"helloWorld": "hello world !"}, *m, "Must dump environment variable in given destination key")
}

func TestCreateEnvDecoratorWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		{
			func() {
				setenv("DECORATORS_ENV_WHATEVER_VALUE", "test")
			},
			`An environment variable suffixed with "DESTKEY" must be defined with "WHATEVER", like DECORATORS_ENV_WHATEVER_DESTKEY`,
		},
		{
			func() {
				setenv("DECORATORS_ENV_WHATEVER_DESTKEY", "test")
			},
			`An environment variable suffixed with "VALUE" must be defined with "WHATEVER", like DECORATORS_ENV_WHATEVER_VALUE`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^DECORATORS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("DECORATORS")

		assert.NoError(t, err, "Must return no errors")

		_, err = createDecorators(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
