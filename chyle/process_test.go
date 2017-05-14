package chyle

import (
	"regexp"
	"testing"

	"github.com/antham/chyle/chyle/config"
	"github.com/antham/chyle/chyle/decorators"
	"github.com/antham/chyle/chyle/extractors"
	"github.com/antham/chyle/chyle/matchers"
	"github.com/antham/chyle/chyle/senders"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	conf := config.CHYLE{}

	p := buildProcess(&conf)

	expected := process{
		&[]matchers.Matcher{},
		&[]extractors.Extracter{},
		&map[string][]decorators.Decorater{},
		&[]senders.Sender{},
	}

	assert.EqualValues(t, expected, *p)
}

func TestBuildProcessWithAFullConfig(t *testing.T) {
	conf := config.CHYLE{}

	conf.FEATURES.MATCHERS.ENABLED = true
	conf.MATCHERS = map[string]string{"TYPE": "merge"}

	conf.FEATURES.EXTRACTORS.ENABLED = true
	conf.EXTRACTORS = map[string]struct {
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

	conf.FEATURES.DECORATORS.ENABLED = true
	conf.FEATURES.DECORATORS.ENABLED = true
	conf.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{
		"TEST": {
			"test",
			"TEST",
		},
	}

	conf.FEATURES.SENDERS.ENABLED = true
	conf.FEATURES.SENDERS.STDOUT = true
	conf.SENDERS.STDOUT.FORMAT = "json"

	p := buildProcess(&conf)

	assert.Len(t, *(p.matchers), 1)
	assert.Len(t, *(p.extractors), 1)
	assert.Len(t, *(p.decorators), 2)
	assert.Len(t, *(p.senders), 1)
}
