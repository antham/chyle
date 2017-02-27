package chyle

import (
	"srcd.works/go-git.v4/plumbing/object"

	"github.com/antham/envh"
)

// process represents all configuration operations defined
// needed to create a changelog
type process struct {
	matchers   *[]matcher
	extractors *[]extracter
	decorators *[]decorater
	senders    *[]sender
}

// buildProcess creates process entity from defined configuration
func buildProcess(config *envh.EnvTree) (*process, error) {
	matchers := &[]matcher{}
	extractors := &[]extracter{}
	decorators := &[]decorater{}
	senders := &[]sender{}

	if subConfig, err := config.FindSubTree("CHYLE", "MATCHERS"); err == nil {
		matchers, err = createMatchers(&subConfig)

		if err != nil {
			return nil, err
		}
	}

	if subConfig, err := config.FindSubTree("CHYLE", "EXTRACTORS"); err == nil {
		extractors, err = createExtractors(&subConfig)

		if err != nil {
			return nil, err
		}
	}

	if subConfig, err := config.FindSubTree("CHYLE", "DECORATORS"); err == nil {
		decorators, err = createDecorators(&subConfig)

		if err != nil {
			return nil, err
		}
	}

	if subConfig, err := config.FindSubTree("CHYLE", "SENDERS"); err == nil {
		senders, err = createSenders(&subConfig)

		if err != nil {
			return nil, err
		}
	}

	return &process{
		matchers,
		extractors,
		decorators,
		senders,
	}, nil
}

// proceed extracts datas from a set of commits
func proceed(process *process, commits *[]object.Commit) error {
	comExt, err := extract(process.extractors, TransformCommitsToMap(filter(process.matchers, commits)))

	if err != nil {

		return err
	}

	comExp, err := decorate(process.decorators, comExt)

	if err != nil {
		return err
	}

	return Send(process.senders, comExp)
}
