package chyle

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestExtract(t *testing.T) {
	extractors := []extracter{
		regexpExtractor{
			"id",
			"serviceId",
			regexp.MustCompile(`(\#\d+)`),
		},
		regexpExtractor{
			"id",
			"booleanValue",
			regexp.MustCompile(`(true|false)`),
		},
		regexpExtractor{
			"id",
			"intValue",
			regexp.MustCompile(` (\d+)`),
		},
		regexpExtractor{
			"id",
			"floatValue",
			regexp.MustCompile(`(\d+\.\d+)`),
		},
		regexpExtractor{
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

	results, err := extract(&extractors, &commitMaps)

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

	setenv("EXTRACTORS_ID_ORIGKEY", "id")
	setenv("EXTRACTORS_ID_DESTKEY", "test")
	setenv("EXTRACTORS_ID_REG", ".*")

	setenv("EXTRACTORS_AUTHORNAME_ORIGKEY", "authorName")
	setenv("EXTRACTORS_AUTHORNAME_DESTKEY", "test2")
	setenv("EXTRACTORS_AUTHORNAME_REG", ".*")

	config, err := envh.NewEnvTree("^EXTRACTORS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("EXTRACTORS")

	assert.NoError(t, err, "Must return no errors")

	e, err := createExtractors(&subConfig)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *e, 2, "Must return 2 extractors")

	expected := map[string]map[string]string{
		"id": map[string]string{
			"index":      "id",
			"identifier": "test",
			"regexp":     ".*",
		},
		"authorName": map[string]string{
			"index":      "authorName",
			"identifier": "test2",
			"regexp":     ".*",
		},
	}

	for i := 0; i < 2; i++ {
		index := (*e)[0].(regexpExtractor).index

		v, ok := expected[index]

		if !ok {
			assert.Fail(t, "Index must exists in expected", "Key must exists")
		}

		assert.Equal(t, (*e)[0].(regexpExtractor).index, v["index"], "Must return first component after extractor variable")
		assert.Equal(t, (*e)[0].(regexpExtractor).identifier, v["identifier"], "Must return second component after extractor variable")
		assert.Equal(t, (*e)[0].(regexpExtractor).re, regexp.MustCompile(v["regexp"]), "Must return value as regexp")
	}
}

func TestCreateExtractorsWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("EXTRACTORS_AUTHORNAME_TEST", "")
			},
			`An environment variable suffixed with "ORIGKEY" must be defined with "AUTHORNAME", like EXTRACTORS_AUTHORNAME_ORIGKEY`,
		},
		g{
			func() {
				setenv("EXTRACTORS_AUTHORNAME_ORIGKEY", "test")
			},
			`An environment variable suffixed with "DESTKEY" must be defined with "AUTHORNAME", like EXTRACTORS_AUTHORNAME_DESTKEY`,
		},
		g{
			func() {
				setenv("EXTRACTORS_AUTHORNAME_ORIGKEY", "test")
				setenv("EXTRACTORS_AUTHORNAME_DESTKEY", "test")
			},
			`An environment variable suffixed with "REG" must be defined with "AUTHORNAME", like EXTRACTORS_AUTHORNAME_REG`,
		},
		g{
			func() {
				setenv("EXTRACTORS_AUTHORNAME_ORIGKEY", "test")
				setenv("EXTRACTORS_AUTHORNAME_DESTKEY", "test")
				setenv("EXTRACTORS_AUTHORNAME_REG", "*")
			},
			`"*" is not a valid regular expression defined for "EXTRACTORS_AUTHORNAME_REG" key`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^EXTRACTORS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("EXTRACTORS")

		assert.NoError(t, err, "Must return no errors")

		_, err = createExtractors(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
