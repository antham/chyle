package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	config, err := envh.NewEnvTree("WHATEVER", "")

	assert.NoError(t, err, "Must produces no errors")

	_, err = buildProcess(&config)

	assert.NoError(t, err, "Must produces no errors")
}

func TestBuildProcessWithWrongEnvironmentsVariables(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		{
			func() {
				setenv("CHYLE_MATCHERS_WHATEVER", "test")
			},
			`a wrong matcher key containing "WHATEVER" was defined`,
		},
		{
			func() {
				setenv("CHYLE_EXTRACTORS_WHATEVER", "test")
			},
			`An environment variable suffixed with "ORIGKEY" must be defined with "WHATEVER", like EXTRACTORS_WHATEVER_ORIGKEY`,
		},
		{
			func() {
				setenv("CHYLE_DECORATORS_WHATEVER", "test")
			},
			`a wrong decorator key containing "WHATEVER" was defined`,
		},
		{
			func() {
				setenv("CHYLE_SENDERS_WHATEVER", "test")
			},
			`a wrong sender key containing "WHATEVER" was defined`,
		},
	}

	for _, test := range tests {
		restoreEnvs()

		test.f()

		config, err := envh.NewEnvTree("CHYLE", "_")

		assert.NoError(t, err, "Must produces no errors")

		_, err = buildProcess(&config)

		assert.EqualError(t, err, test.e, "Must returns an error")
	}
}
