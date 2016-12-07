package chyle

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStdoutSender(t *testing.T) {
	s, err := NewStdoutSender("json")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, StdoutSender{"json", os.Stdout}, s, "Must return a stdout sender")
}

func TestNewStdoutSenderWithAnUnexistingFormat(t *testing.T) {
	s, err := NewStdoutSender("whatever")

	assert.EqualError(t, err, "\"whatever\" format does not exist", "Must return an error")
	assert.Equal(t, StdoutSender{}, s, "Must return a stdout sender")
}

func TestStdoutSenderWithJson(t *testing.T) {
	buf := &bytes.Buffer{}

	s := StdoutSender{"json", buf}
	err := s.Send(&[]map[string]interface{}{
		map[string]interface{}{
			"id":   1,
			"test": "test",
		},
		map[string]interface{}{
			"id":   2,
			"test": "test",
		},
	})

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `[{"id":1,"test":"test"},{"id":2,"test":"test"}]`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations  as json")
}

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
	r, err := CreateSenders(map[string]interface{}{
		"stdout": map[string]string{
			"format": "json",
		},
	})

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *r, 1, "Must return 1 expander")
}

func TestCreateSendersWithErrors(t *testing.T) {
	type g struct {
		s map[string]interface{}
		e string
	}

	datas := []g{
		g{
			map[string]interface{}{"whatever": map[string]string{"test": "test"}},
			`"whatever" is not a valid sender structure`,
		},
		g{
			map[string]interface{}{"stdout": map[string]string{"test": "test"}},
			`"format" key must be defined`,
		},
		g{
			map[string]interface{}{"stdout": map[string]string{"format": "test"}},
			`"test" format does not exist`,
		},
	}

	for _, d := range datas {
		_, err := CreateSenders(d.s)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, d.e, "Must match error string")
	}
}
