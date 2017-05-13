package chyle

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	p := buildProcess()

	expected := process{
		&[]matcher{},
		&[]extracter{},
		&map[string][]decorater{},
		&[]sender{},
	}

	assert.EqualValues(t, expected, *p)
}

func TestBuildProcessWithAFullConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	chyleConfig.FEATURES.HASMATCHERS = true
	chyleConfig.MATCHERS = map[string]string{"TYPE": "merge"}
	chyleConfig.FEATURES.HASEXTRACTORS = true
	chyleConfig.EXTRACTORS = map[string]struct {
		ORIGKEY string
		DESTKEY string
		REG     *regexp.Regexp
	}{
		"TEST": {
			"TEST",
			"test",
			regexp.MustCompile(".*"),
		},
	}
	chyleConfig.FEATURES.HASDECORATORS = true
	chyleConfig.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"TEST": {
			"test",
			"TEST",
		},
	}
	chyleConfig.FEATURES.HASSTDOUTSENDER = true
	chyleConfig.SENDERS.STDOUT.FORMAT = "json"

	p := buildProcess()

	expected := process{
		&[]matcher{
			mergeCommitMatcher{},
		},
		&[]extracter{
			regexpExtractor{
				index:      "TEST",
				identifier: "test",
				re:         regexp.MustCompile(".*"),
			},
		},
		&map[string][]decorater{
			"datas":     {},
			"metadatas": {},
		},
		&[]sender{},
	}

	assert.Equal(t, expected, *p)
}
