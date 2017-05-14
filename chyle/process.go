package chyle

import (
	"github.com/antham/chyle/chyle/decorators"
	"github.com/antham/chyle/chyle/extractors"
	"github.com/antham/chyle/chyle/matchers"
	"github.com/antham/chyle/chyle/senders"

	"srcd.works/go-git.v4/plumbing/object"
)

// process represents all configuration operations defined
// needed to create a changelog
type process struct {
	matchers   *[]matchers.Matcher
	extractors *[]extractors.Extracter
	decorators *map[string][]decorators.Decorater
	senders    *[]senders.Sender
}

// buildProcess creates process entity from defined configuration
func buildProcess() *process {
	p := &process{
		&[]matchers.Matcher{},
		&[]extractors.Extracter{},
		&map[string][]decorators.Decorater{},
		&[]senders.Sender{},
	}

	if chyleConfig.FEATURES.MATCHERS.ENABLED {
		p.matchers = matchers.CreateMatchers(chyleConfig.MATCHERS)
	}

	if chyleConfig.FEATURES.EXTRACTORS.ENABLED {
		p.extractors = extractors.CreateExtractors(chyleConfig.EXTRACTORS)
	}

	if chyleConfig.FEATURES.DECORATORS.ENABLED {
		p.decorators = decorators.CreateDecorators(chyleConfig.FEATURES.DECORATORS, chyleConfig.DECORATORS)
	}

	if chyleConfig.FEATURES.SENDERS.ENABLED {
		p.senders = senders.CreateSenders(chyleConfig.FEATURES.SENDERS, chyleConfig.SENDERS)
	}

	return p
}

// proceed extracts datas from a set of commits
func proceed(process *process, commits *[]object.Commit) error {
	changelog := extractors.Extract(process.extractors, matchers.Filter(process.matchers, commits))

	changelog, err := decorators.Decorate(process.decorators, changelog)

	if err != nil {
		return err
	}

	return senders.Send(process.senders, changelog)
}
