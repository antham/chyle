package chyle

import (
	"regexp"
	"testing"

	"github.com/antham/chyle/chyle/decorators"
	"github.com/antham/chyle/chyle/extractors"
	"github.com/antham/chyle/chyle/matchers"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	p := buildProcess()

	expected := process{
		&[]matchers.Matcher{},
		&[]extractors.Extracter{},
		&map[string][]decorators.Decorater{},
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
	chyleConfig.FEATURES.HASENVDECORATOR = true
	chyleConfig.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"TEST": {
			"test",
			"TEST",
		},
	}

	chyleConfig.FEATURES.HASSENDERS = true
	chyleConfig.FEATURES.HASSTDOUTSENDER = true
	chyleConfig.SENDERS.STDOUT.FORMAT = "json"

	p := buildProcess()

	assert.Len(t, *(p.matchers), 1)
	assert.Len(t, *(p.extractors), 1)
	assert.Len(t, *(p.decorators), 2)
	assert.Len(t, *(p.senders), 1)
}
