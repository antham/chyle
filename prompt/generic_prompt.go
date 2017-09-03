package prompt

import (
	"fmt"
	"strings"

	"github.com/antham/chyle/prompt/internal/counter"
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

func parseEnv(env string, store *Store) func(value string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("No value given")
		}

		(*store)[env] = value

		return nil
	}
}

func parseEnvWithCounter(env string, counter *counter.Counter, store *Store) func(value string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("No value given")
		}

		(*store)[strings.Replace(env, "*", counter.Get(), -1)] = value

		return nil
	}
}

func parseEnvWithCounterAndIncrement(env string, counter *counter.Counter, store *Store) func(value string) error {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("No value given")
		}

		(*store)[strings.Replace(env, "*", counter.Get(), -1)] = value

		counter.Increment()

		return nil
	}
}

type envConfig struct {
	ID           string
	nextID       string
	env          string
	promptString string
}

func newEnvPrompt(config envConfig, store *Store) strumt.Prompter {
	return &genericPrompt{
		config.ID,
		config.promptString,
		func(string) string { return config.nextID },
		func(error) string { return config.ID },
		parseEnv(config.env, store),
	}
}

func newGroupEnvPromptWithCounter(configs []envConfig, store *Store) []strumt.Prompter {
	results := []strumt.Prompter{}
	c := &counter.Counter{}

	for i, config := range configs {
		f := parseEnvWithCounter(config.env, c, store)

		if i == len(configs)-1 {
			f = parseEnvWithCounterAndIncrement(config.env, c, store)
		}

		p := genericPrompt{
			config.ID,
			config.promptString,
			func(nextID string) func(string) string { return func(string) string { return nextID } }(config.nextID),
			func(ID string) func(error) string { return func(error) string { return ID } }(config.ID),
			f,
		}

		results = append(results, &p)
	}

	return results
}
	}
}
