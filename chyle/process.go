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

	if chyleConfig.FEATURES.HASMATCHERS {
		p.matchers = matchers.CreateMatchers(chyleConfig.MATCHERS)
	}

	if chyleConfig.FEATURES.HASEXTRACTORS {
		p.extractors = extractors.CreateExtractors(chyleConfig.EXTRACTORS)
	}

	if chyleConfig.FEATURES.HASDECORATORS {
		p.decorators = decorators.CreateDecorators(map[string]bool{"jiraIssueDecorator": chyleConfig.FEATURES.HASJIRAISSUEDECORATOR, "githubIssueDecorator": chyleConfig.FEATURES.HASGITHUBISSUEDECORATOR, "envDecorator": chyleConfig.FEATURES.HASENVDECORATOR}, chyleConfig.DECORATORS)
	}

	if chyleConfig.FEATURES.HASSENDERS {
		p.senders = senders.CreateSenders(map[string]bool{"githubReleaseSender": chyleConfig.FEATURES.HASGITHUBRELEASESENDER, "stdoutSender": chyleConfig.FEATURES.HASSTDOUTSENDER}, chyleConfig.SENDERS)
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
