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
