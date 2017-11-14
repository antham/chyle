package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

var mainMenu = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"mainMenu",
		addQuitChoice(
			[]builder.SwitchConfig{
				{"1", "Add a matcher", "matcherChoice"},
				{"2", "Add an extractor", "extractorOrigKey"},
				{"3", "Add a decorator", "decoratorChoice"},
				{"4", "Add a sender", "senderChoice"},
			},
		),
	),
}

func newMainMenu() []strumt.Prompter {
	return mainMenu
}
