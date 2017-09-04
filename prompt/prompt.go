package prompt

import (
	"github.com/antham/strumt"
)

// Prompts held prompts
type Prompts struct {
	prompts strumt.Prompts
}

// New creates a new prompt chain
func New() Prompts {
	return Prompts{strumt.NewPrompts()}
}

// Run starts a prompt session
func (p *Prompts) Run() (*Store, error) {
	store := &Store{}
	prompts := mergePrompters(
		newMainMenu(),
		newMandatoryOption(store),
		newMatchers(store),
		newExtractors(store),
		newDecorators(store),
		newSenders(store),
	)

	for _, item := range prompts {
		switch prompt := item.(type) {
		case strumt.LinePrompter:
			p.prompts.AddLinePrompter(prompt)
		case strumt.MultilinePrompter:
			p.prompts.AddMultilinePrompter(prompt)
		}
	}

	p.prompts.SetFirst("referenceFrom")
	p.prompts.Run()

	return store, nil
}
