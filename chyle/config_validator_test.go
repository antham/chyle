package chyle

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"
)

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
