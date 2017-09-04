package prompt

import (
	"github.com/antham/strumt"
)

func newMandatoryOption(store *Store) []strumt.Prompter {
	return newEnvPrompts(mandatoryOption, store)
}

var mandatoryOption = []envConfig{
	envConfig{"referenceFrom", "referenceTo", "CHYLE_GIT_REFERENCE_FROM", "Enter a git commit ID that start your range"},
	envConfig{"referenceTo", "gitPath", "CHYLE_GIT_REFERENCE_TO", "Enter a git commit ID that end your range"},
	envConfig{"gitPath", "mainMenu", "CHYLE_GIT_REPOSITORY_PATH", "Enter your git path repository"},
}
