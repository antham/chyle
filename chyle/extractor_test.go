package chyle

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestExtract(t *testing.T) {
	extractors := []Extracter{
		RegexpExtracter{
			"id",
			"serviceId",
			regexp.MustCompile(`(\#\d+)`),
		},
		RegexpExtracter{
			"id",
			"booleanValue",
			regexp.MustCompile(`(true|false)`),
		},
		RegexpExtracter{
			"id",
			"intValue",
			regexp.MustCompile(` (\d+)`),
		},
		RegexpExtracter{
			"id",
			"floatValue",
			regexp.MustCompile(`(\d+\.\d+)`),
		},
		RegexpExtracter{
			"secondIdentifier",
			"secondServiceId",
			regexp.MustCompile(`(#\d+)`),
		},
	}

	commitMaps := []map[string]interface{}{
		map[string]interface{}{
			"id":               "Whatever #30 whatever true 12345 whatever 12345.12",
			"secondIdentifier": "test #12345",
		},
		map[string]interface{}{
			"id":               "Whatever #40 whatever false whatever 78910 whatever 78910.12",
			"secondIdentifier": "test #45678",
		},
		map[string]interface{}{
			"id": "Whatever whatever whatever",
		},
	}

	results, err := Extract(&extractors, &commitMaps)

	expected := []map[string]interface{}{
		map[string]interface{}{
			"id":               "Whatever #30 whatever true 12345 whatever 12345.12",
			"secondIdentifier": "test #12345",
			"serviceId":        "#30",
			"secondServiceId":  "#12345",
			"booleanValue":     true,
			"intValue":         int64(12345),
			"floatValue":       12345.12,
		},
		map[string]interface{}{
			"id":               "Whatever #40 whatever false whatever 78910 whatever 78910.12",
			"secondIdentifier": "test #45678",
			"serviceId":        "#40",
			"secondServiceId":  "#45678",
			"booleanValue":     false,
			"intValue":         int64(78910),
			"floatValue":       78910.12,
		},
		map[string]interface{}{
			"id":           "Whatever whatever whatever",
			"serviceId":    "",
			"booleanValue": "",
			"intValue":     "",
			"floatValue":   "",
		},
	}

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, expected, *results, "Must return extracted datas with old one")
}

func TestCreateExtractors(t *testing.T) {
	restoreEnvs()

	setenv("EXTRACTORS_ID_TEST", ".*")
	setenv("EXTRACTORS_AUTHORNAME_TEST2", ".*")

	config, err := envh.NewEnvTree("^EXTRACTORS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("EXTRACTORS")

	assert.NoError(t, err, "Must return no errors")

	e, err := CreateExtractors(&subConfig)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *e, 2, "Must return 2 extractors")

	assert.Equal(t, (*e)[0].(RegexpExtracter).index, "id", "Must return first component after extractor variable")
	assert.Equal(t, (*e)[0].(RegexpExtracter).identifier, "test", "Must return second component after extractor variable")
	assert.Equal(t, (*e)[0].(RegexpExtracter).re, regexp.MustCompile(".*"), "Must return value as regexp")

	assert.Equal(t, (*e)[1].(RegexpExtracter).index, "authorname", "Must return first component after extractor variable")
	assert.Equal(t, (*e)[1].(RegexpExtracter).identifier, "test2", "Must return second component after extractor variable")
	assert.Equal(t, (*e)[1].(RegexpExtracter).re, regexp.MustCompile(".*"), "Must return value as regexp")
}

func TestCreateExtractorsWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("EXTRACTORS_AUTHORNAME_TEST", "*")
			},
			`"*" is not a valid regular expression defined for "AUTHORNAME" in "EXTRACTORS" config`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^EXTRACTORS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("EXTRACTORS")

		assert.NoError(t, err, "Must return no errors")

		_, err = CreateExtractors(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
