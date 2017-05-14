package chyle

import (
	"regexp"
	"testing"

	"github.com/antham/chyle/chyle/decorators"
	"github.com/antham/chyle/chyle/extractors"
	"github.com/antham/chyle/chyle/matchers"
	"github.com/antham/chyle/chyle/senders"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	p := buildProcess()

	expected := process{
		&[]matchers.Matcher{},
		&[]extractors.Extracter{},
		&map[string][]decorators.Decorater{},
		&[]senders.Sender{},
	}

	assert.EqualValues(t, expected, *p)
}

func TestBuildProcessWithAFullConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	chyleConfig.FEATURES.MATCHERS.ENABLED = true
	chyleConfig.MATCHERS = map[string]string{"TYPE": "merge"}

	chyleConfig.FEATURES.EXTRACTORS.ENABLED = true
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

	chyleConfig.FEATURES.DECORATORS.ENABLED = true
	chyleConfig.FEATURES.DECORATORS.ENABLED = true
	chyleConfig.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"TEST": {
			"test",
			"TEST",
		},
	}

	chyleConfig.FEATURES.SENDERS.ENABLED = true
	chyleConfig.FEATURES.SENDERS.STDOUT = true
	chyleConfig.SENDERS.STDOUT.FORMAT = "json"

	p := buildProcess()

	assert.Len(t, *(p.matchers), 1)
	assert.Len(t, *(p.extractors), 1)
	assert.Len(t, *(p.decorators), 2)
	assert.Len(t, *(p.senders), 1)
}
