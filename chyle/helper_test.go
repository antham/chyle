package chyle

import (
	"bytes"
	"log"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"
)

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

func TestExtractStringConfig(t *testing.T) {
	restoreEnvs()

	setenv("CHYLE_TEST1", "test1")

	config, err := envh.NewEnvTree("CHYLE", "_")

	if err != nil {
		logrus.Fatal(err)
	}

	var test1 string
	var test2 string

	err = extractStringConfig(
		&config,
		[]strConfigMapping{
			strConfigMapping{
				[]string{"CHYLE", "TEST1"},
				&test1,
				true,
			},
			strConfigMapping{
				[]string{"CHYLE", "TEST2"},
				&test2,
				false,
			},
		},
		[]string{""},
	)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test1", test1, "Must return test1")
	assert.Equal(t, "", test2, "Must return nothing, variable is not defined")
}

func TestExtractBoolConfig(t *testing.T) {
	restoreEnvs()

	setenv("CHYLE_TEST1", "true")

	config, err := envh.NewEnvTree("CHYLE", "_")

	if err != nil {
		logrus.Fatal(err)
	}

	var test1 bool
	var test2 bool

	err = extractBoolConfig(
		&config,
		[]boolConfigMapping{
			boolConfigMapping{
				[]string{"CHYLE", "TEST1"},
				&test1,
				true,
			},
			boolConfigMapping{
				[]string{"CHYLE", "TEST2"},
				&test2,
				false,
			},
		},
		[]string{""},
	)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, true, test1, "Must return false")
	assert.Equal(t, false, test2, "Must return default value cause variable is not defined")
}
