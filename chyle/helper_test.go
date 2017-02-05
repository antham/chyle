package chyle

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopulateTemplate(t *testing.T) {
	r, err := populateTemplate("test", "{{.test}}", map[string]string{"test": "Hello world !"})

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "Hello world !", r, "Must interpolate data in template")
}

func TestPopulateTemplateWithABadTemplate(t *testing.T) {
	_, err := populateTemplate("test", "{{.test", map[string]string{"test": "Hello world !"})

	assert.EqualError(t, err, "check your template is well-formed : template: test:1: unclosed action", "Must return an error when template is not well-formed")
}

func TestDebug(t *testing.T) {
	b := []byte{}

	buffer := bytes.NewBuffer(b)

	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

	EnableDebugging = true

	debug("test : %s", "output")

	actual, err := buffer.ReadString('\n')

	assert.NoError(t, err, "Must return no errors")
	assert.Regexp(t, `CHYLE - \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} test : output\n`, actual, "Must output given format with argument when debug is enabled")
}

func TestDebugWithDebugDisabled(t *testing.T) {
	b := []byte{}

	buffer := bytes.NewBuffer(b)

	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

	EnableDebugging = false

	debug("test : %s", "output")

	_, err := buffer.ReadString('\n')

	assert.EqualError(t, err, "EOF", "Must return EOF error")
}
