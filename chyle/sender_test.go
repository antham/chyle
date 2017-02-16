package chyle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestSend(t *testing.T) {
	buf := &bytes.Buffer{}

	s := StdoutSender{"json", buf}
	datas := &[]map[string]interface{}{
		map[string]interface{}{
			"id":   1,
			"test": "test",
		},
		map[string]interface{}{
			"id":   2,
			"test": "test",
		},
	}

	err := Send(&[]Sender{s}, datas)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `[{"id":1,"test":"test"},{"id":2,"test":"test"}]`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations  as json")
}

func TestCreateSenders(t *testing.T) {
	restoreEnvs()
	setenv("SENDERS_STDOUT_FORMAT", "json")

	config, err := envh.NewEnvTree("^SENDERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("SENDERS")

	assert.NoError(t, err, "Must return no errors")

	r, err := CreateSenders(&subConfig)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *r, 1, "Must return 1 decorator")
}

func TestCreateSendersWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("SENDERS_WHATEVER", "test")
			},
			`a wrong sender key containing "WHATEVER" was defined`,
		},
		g{
			func() {
				setenv("SENDERS_STDOUT_FORMAT", "test")
			},
			`"test" format does not exist`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^SENDERS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("SENDERS")

		assert.NoError(t, err, "Must return no errors")

		_, err = CreateSenders(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
