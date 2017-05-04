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
	chyleConfig.EXTRACTORS = map[string]map[string]string{"TEST": {"TEST": "test"}}
	chyleConfig.FEATURES.HASDECORATORS = true
	chyleConfig.DECORATORS.ENV = map[string]map[string]string{"TEST": {"TEST": "test"}}
	chyleConfig.FEATURES.HASSTDOUTSENDER = true
	chyleConfig.SENDERS.STDOUT.FORMAT = "json"

	p := buildProcess()

	expected := process{
		&[]matcher{
			mergeCommitMatcher{},
		},
		&[]extracter{
			regexpExtractor{
				re: regexp.MustCompile(""),
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
