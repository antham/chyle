package prompt

import (
	"fmt"

	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newMatchers(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		matcherChoice,
		builder.NewEnvPrompts(matcher, store),
	)
}

var matcherChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"matcherChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				{
					Choice:       "1",
					PromptString: "Add a type matcher",
					NextPromptID: "matcherType",
				},
				{
					Choice:       "2",
					PromptString: "Add a message matcher",
					NextPromptID: "matcherMessage",
				},
				{
					Choice:       "3",
					PromptString: "Add a committer matcher",
					NextPromptID: "matcherCommitter",
				},
				{
					Choice:       "4",
					PromptString: "Add an author matcher",
					NextPromptID: "matcherAuthor",
				},
			},
		),
	),
}

var matcher = []builder.EnvConfig{
	{
		ID:           "matcherType",
		NextID:       "matcherChoice",
		Env:          "CHYLE_MATCHERS_TYPE",
		PromptString: "Enter a matcher type (regular or merge)",
		Validator:    validateMatcherType,
	},
	{
		ID:           "matcherMessage",
		NextID:       "matcherChoice",
		Env:          "CHYLE_MATCHERS_MESSAGE",
		PromptString: "Enter a regexp to match commit message",
		Validator:    validateRegexp,
	},
	{
		ID:           "matcherCommitter",
		NextID:       "matcherChoice",
		Env:          "CHYLE_MATCHERS_COMMITTER",
		PromptString: "Enter a regexp to match git committer",
		Validator:    validateRegexp,
	},
	{
		ID:           "matcherAuthor",
		NextID:       "matcherChoice",
		Env:          "CHYLE_MATCHERS_AUTHOR",
		PromptString: "Enter a regexp to match git author",
		Validator:    validateRegexp,
	},
}

func validateMatcherType(value string) error {
	if value != "regular" && value != "merge" {
		return fmt.Errorf(`Must be "regular" or "merge"`)
	}

	return nil
}
