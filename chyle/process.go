package chyle

import (
	pkgDecorators "github.com/antham/chyle/chyle/decorators"
	pkgExtractors "github.com/antham/chyle/chyle/extractors"
	pkgMatchers "github.com/antham/chyle/chyle/matchers"
	pkgSenders "github.com/antham/chyle/chyle/senders"

	"srcd.works/go-git.v4/plumbing/object"
)

// process represents all configuration operations defined
// needed to create a changelog
type process struct {
	matchers   *[]pkgMatchers.Matcher
	extractors *[]pkgExtractors.Extracter
	decorators *map[string][]pkgDecorators.Decorater
	senders    *[]pkgSenders.Sender
}

// buildProcess creates process entity from defined configuration
func buildProcess() *process {
	matchers := &[]pkgMatchers.Matcher{}
	extractors := &[]pkgExtractors.Extracter{}
	decorators := &map[string][]pkgDecorators.Decorater{}
	senders := &[]pkgSenders.Sender{}

	if chyleConfig.FEATURES.HASMATCHERS {
		matchers = pkgMatchers.CreateMatchers(chyleConfig.MATCHERS)
	}

	if chyleConfig.FEATURES.HASEXTRACTORS {
		extractors = pkgExtractors.CreateExtractors(chyleConfig.EXTRACTORS)
	}

	if chyleConfig.FEATURES.HASDECORATORS {
		decorators = pkgDecorators.CreateDecorators(map[string]bool{"jiraIssueDecorator": chyleConfig.FEATURES.HASJIRAISSUEDECORATOR, "githubIssueDecorator": chyleConfig.FEATURES.HASGITHUBISSUEDECORATOR, "envDecorator": chyleConfig.FEATURES.HASENVDECORATOR}, chyleConfig.DECORATORS)
	}

	if chyleConfig.FEATURES.HASSENDERS {
		senders = pkgSenders.CreateSenders(map[string]bool{"githubReleaseSender": chyleConfig.FEATURES.HASGITHUBRELEASESENDER, "stdoutSender": chyleConfig.FEATURES.HASSTDOUTSENDER}, chyleConfig.SENDERS)
	}

	return &process{
		matchers,
		extractors,
		decorators,
		senders,
	}
}

// proceed extracts datas from a set of commits
func proceed(process *process, commits *[]object.Commit) error {
	changelog := pkgExtractors.Extract(process.extractors, pkgMatchers.Filter(process.matchers, commits))

	changelog, err := pkgDecorators.Decorate(process.decorators, changelog)

	if err != nil {
		return err
	}

	return pkgSenders.Send(process.senders, changelog)
}
