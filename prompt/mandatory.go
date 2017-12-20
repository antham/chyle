package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newMandatoryOption(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(mandatoryOption, store)
}

var mandatoryOption = []builder.EnvConfig{
	{
		ID:           "referenceFrom",
		NextID:       "referenceTo",
		Env:          "CHYLE_GIT_REFERENCE_FROM",
		PromptString: "Enter a git commit ID that start your range",
		Validator:    noOpValidator,
	},
	{
		ID:           "referenceTo",
		NextID:       "gitPath",
		Env:          "CHYLE_GIT_REFERENCE_TO",
		PromptString: "Enter a git commit ID that end your range",
		Validator:    noOpValidator,
	},
	{
		ID:           "gitPath",
		NextID:       "mainMenu",
		Env:          "CHYLE_GIT_REPOSITORY_PATH",
		PromptString: "Enter your git path repository",
		Validator:    noOpValidator,
	},
}
