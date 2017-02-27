package chyle

import (
	"os"
	"testing"

	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"
)

func TestBuildStdoutSender(t *testing.T) {
	restoreEnvs()
	setenv("SENDERS_STDOUT_FORMAT", "json")

	config, err := envh.NewEnvTree("^SENDERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("SENDERS", "STDOUT")

	assert.NoError(t, err, "Must return no errors")

	s, err := buildStdoutSender(&subConfig)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, jSONStdoutSender{os.Stdout}, s, "Must return a json stdout sender")
}

func TestBuildStdoutSenderWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("SENDERS_STDOUT_FORMAT", "test")
			},
			`"test" format does not exist`,
		},
		g{
			func() {
				setenv("SENDERS_STDOUT_FORMAT", "template")
			},
			`"SENDERS_STDOUT_TEMPLATE" must be defined when "template" format is defined`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^SENDERS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("SENDERS", "STDOUT")

		assert.NoError(t, err, "Must return no errors")

		_, err = buildStdoutSender(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
