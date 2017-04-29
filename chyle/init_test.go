package chyle

import (
// "bytes"
// "log"
// "testing"

// "github.com/stretchr/testify/assert"
)

// func TestDebug(t *testing.T) {
// 	b := []byte{}

// 	buffer := bytes.NewBuffer(b)

// 	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

// 	EnableDebugging = true

// 	debug("test : %s", "output")

// 	actual, err := buffer.ReadString('\n')

// 	assert.NoError(t, err, "Must return no errors")
// 	assert.Regexp(t, `CHYLE - \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} test : output\n`, actual, "Must output given format with argument when debug is enabled")
// }

// func TestDebugWithDebugDisabled(t *testing.T) {
// 	b := []byte{}

// 	buffer := bytes.NewBuffer(b)

// 	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

// 	EnableDebugging = false

// 	debug("test : %s", "output")

// 	_, err := buffer.ReadString('\n')

// 	assert.EqualError(t, err, "EOF", "Must return EOF error")
// }
