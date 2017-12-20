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
				{
					Choice:       "1",
					PromptString: "Add a matcher",
					NextPromptID: "matcherChoice",
				},
				{
					Choice:       "2",
					PromptString: "Add an extractor",
					NextPromptID: "extractorOrigKey",
				},
				{
					Choice:       "3",
					PromptString: "Add a decorator",
					NextPromptID: "decoratorChoice",
				},
				{
					Choice:       "4",
					PromptString: "Add a sender",
					NextPromptID: "senderChoice",
				},
			},
		),
	),
}

func newMainMenu() []strumt.Prompter {
	return mainMenu
}
