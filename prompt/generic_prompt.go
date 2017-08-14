package prompt

import (
	"fmt"
	"os"

	"github.com/antham/strumt"
)

type genericPrompt struct {
	iD           string
	promptString string
	onSuccess    func(string) string
	onError      func(error) string
	parse        func(string) error
}

func (g *genericPrompt) ID() string {
	return g.iD
}

func (g *genericPrompt) PromptString() string {
	return g.promptString
}

func (g *genericPrompt) Parse(value string) error {
	return g.parse(value)
}

func (g *genericPrompt) NextOnSuccess(value string) string {
	return g.onSuccess(value)
}

func (g *genericPrompt) NextOnError(err error) string {
	return g.onError(err)
}

func parseEnv(env string) func(value string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("No value given")
		}

		return os.Setenv(env, value)
	}
}

type envConfig struct {
	ID           string
	nextID       string
	env          string
	promptString string
}

func newEnvPrompt(config envConfig) strumt.Prompter {
	return &genericPrompt{
		config.ID,
		config.promptString,
		func(string) string { return config.nextID },
		func(error) string { return config.ID },
		parseEnv(config.env),
	}
}
