package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newMandatoryOption(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(mandatoryOption, store)
}

var mandatoryOption = []builder.EnvConfig{
	builder.EnvConfig{"referenceFrom", "referenceTo", "CHYLE_GIT_REFERENCE_FROM", "Enter a git commit ID that start your range"},
	builder.EnvConfig{"referenceTo", "gitPath", "CHYLE_GIT_REFERENCE_TO", "Enter a git commit ID that end your range"},
	builder.EnvConfig{"gitPath", "mainMenu", "CHYLE_GIT_REPOSITORY_PATH", "Enter your git path repository"},
}
