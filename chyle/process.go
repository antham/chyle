package chyle

import (
	"srcd.works/go-git.v4/plumbing/object"
)

// process represents all configuration operations defined
// needed to create a changelog
type process struct {
	matchers   *[]matcher
	extractors *[]extracter
	decorators *map[string][]decorater
	senders    *[]sender
}

// buildProcess creates process entity from defined configuration
func buildProcess() (*process, error) {
	matchers := &[]matcher{}
	extractors := &[]extracter{}
	decorators := &map[string][]decorater{}
	senders := &[]sender{}

	if chyleConfig.FEATURES.HASMATCHERS {
		matchers = createMatchers()
	}

	if chyleConfig.FEATURES.HASEXTRACTORS {
		extractors = createExtractors()
	}

	if chyleConfig.FEATURES.HASDECORATORS {
		decorators = createDecorators()
	}

	if chyleConfig.FEATURES.HASSENDERS {
		senders = createSenders()
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
	changelog, err := extract(process.extractors, TransformCommitsToMap(filter(process.matchers, commits)))

	if err != nil {

		return err
	}

	changelog, err = decorate(process.decorators, changelog)

	if err != nil {
		return err
	}

	return Send(process.senders, changelog)
}
