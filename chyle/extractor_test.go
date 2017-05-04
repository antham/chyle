package chyle

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{
			"id":               "Whatever #30 whatever true 12345 whatever 12345.12",
			"secondIdentifier": "test #12345",
		},
		{
			"id":               "Whatever #40 whatever false whatever 78910 whatever 78910.12",
			"secondIdentifier": "test #45678",
		},
		{
			"id": "Whatever whatever whatever",
		},
	}

	results := extract(&extractors, &commitMaps)

	expected := Changelog{
		Datas: []map[string]interface{}{
			{
				"id":               "Whatever #30 whatever true 12345 whatever 12345.12",
				"secondIdentifier": "test #12345",
				"serviceId":        "#30",
				"secondServiceId":  "#12345",
				"booleanValue":     true,
				"intValue":         int64(12345),
				"floatValue":       12345.12,
			},
			{
				"id":               "Whatever #40 whatever false whatever 78910 whatever 78910.12",
				"secondIdentifier": "test #45678",
				"serviceId":        "#40",
				"secondServiceId":  "#45678",
				"booleanValue":     false,
				"intValue":         int64(78910),
				"floatValue":       78910.12,
			},
			{
				"id":           "Whatever whatever whatever",
				"serviceId":    "",
				"booleanValue": "",
				"intValue":     "",
				"floatValue":   "",
			},
		},
		Metadatas: map[string]interface{}{},
	}

	assert.Equal(t, expected, *results, "Must return extracted datas with old one")
}

func TestCreateExtractors(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASEXTRACTORS = true
	chyleConfig.EXTRACTORS = map[string]map[string]string{}

	chyleConfig.EXTRACTORS["ID"] = map[string]string{}
	chyleConfig.EXTRACTORS["ID"]["ORIGKEY"] = "id"
	chyleConfig.EXTRACTORS["ID"]["DESTKEY"] = "test"
	chyleConfig.EXTRACTORS["ID"]["REG"] = ".*"

	chyleConfig.EXTRACTORS["AUTHORNAME"] = map[string]string{}
	chyleConfig.EXTRACTORS["AUTHORNAME"]["ORIGKEY"] = "authorName"
	chyleConfig.EXTRACTORS["AUTHORNAME"]["DESTKEY"] = "test2"
	chyleConfig.EXTRACTORS["AUTHORNAME"]["REG"] = ".*"

	e := createExtractors()

	assert.Len(t, *e, 2, "Must return 2 extractors")

	expected := map[string]map[string]string{
		"id": {
			"index":      "id",
			"identifier": "test",
			"regexp":     ".*",
		},
		"authorName": {
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
