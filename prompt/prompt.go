package prompt

import (
	"github.com/antham/strumt"
)

var mandatoryOption = []strumt.Prompter{
	newEnvPrompt(envConfig{"referenceFrom", "referenceTo", "CHYLE_GIT_REPOSITORY_PATH", "Enter a git commit ID that start your range"}),
	newEnvPrompt(envConfig{"referenceTo", "gitPath", "CHYLE_GIT_REFERENCE_FROM", "Enter a git commit ID that end your range"}),
	newEnvPrompt(envConfig{"gitPath", "matchChoice", "CHYLE_GIT_REFERENCE_TO", "Enter your git path repository"}),
}

// Run starts a prompt session
func Run() ([]string, error) {
	envs := []string{}

	p := strumt.NewPrompts()

	prompts := append([]strumt.Prompter{}, mandatoryOption...)

	for _, item := range prompts {
		switch prompt := item.(type) {
		case strumt.LinePrompter:
			p.AddLinePrompter(prompt)
		case strumt.MultilinePrompter:
			p.AddMultilinePrompter(prompt)
		}
	}

	p.SetFirst("referenceFrom")
	p.Run()

	return envs, nil
}
