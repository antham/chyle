package chyle

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONStdoutSender(t *testing.T) {
	buf := &bytes.Buffer{}

	s := jSONStdoutSender{buf}
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
	assert.Equal(t, `[{"id":1,"test":"test"},{"id":2,"test":"test"}]`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations as json")
}

func TestTemplateStdoutSender(t *testing.T) {
	buf := &bytes.Buffer{}

	s := templateStdoutSender{"{{ range $key, $value := . }}{{$value.id}} : {{$value.test}} | {{ end }}", buf}
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
	assert.Equal(t, `1 : test | 2 : test | `, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations shaped on given template")
}
