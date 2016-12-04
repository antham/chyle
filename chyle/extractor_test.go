package chyle

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
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
	e, err := CreateExtractors(
		map[string]interface{}{
			"id":         map[string]string{"test": ".*"},
			"authorName": map[string]string{"test2": ".*"},
		},
	)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *e, 2, "Must return 2 extractors")
}

func TestCreateExtractorsWithErrors(t *testing.T) {
	type g struct {
		s map[string]interface{}
		e string
	}

	datas := []g{
		g{
			map[string]interface{}{"id": map[string]string{"test": "**"}},
			`"test" doesn't contain a valid regular expression`,
		},
	}

	for _, d := range datas {
		_, err := CreateExtractors(d.s)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, d.e, "Must match error string")
	}
}
