package prompt

import (
	"github.com/antham/strumt"
)

var mainMenu = []strumt.Prompter{
	&switchPrompt{
		"mainMenu",
		addQuitChoice(
			[]switchChoice{
				switchChoice{"1", "Add a matcher", "matcherChoice"},
				switchChoice{"2", "Add an extractor", "extractorOrigKey"},
				switchChoice{"3", "Add a decorator", "decoratorChoice"},
				switchChoice{"4", "Add a sender", "senderChoice"},
			},
		),
	},
}

func newMainMenu() []strumt.Prompter {
	return mainMenu
}
