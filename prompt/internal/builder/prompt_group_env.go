package builder

import (
	"fmt"
	"strings"

	"github.com/antham/chyle/prompt/internal/counter"
	"github.com/antham/strumt"
)

// NewGroupEnvPromptWithCounter gives the ability to create several group of related environment variable, a common prefix provided as a number from an internal counter tied variable together. For instance in variables environments TEST_*_KEY and TEST_*_VALUE, * is replaced with a number, it becomes TEST_0_KEY and TEST_0_VALUE another call would give TEST_1_VALUE and TEST_1_KEY
func NewGroupEnvPromptWithCounter(configs []EnvConfig, store *Store) []strumt.Prompter {
	results := []strumt.Prompter{}
	c := &counter.Counter{}

	for i, config := range configs {
		f := parseEnvWithCounter(config.Env, c, store)

		if i == len(configs)-1 {
			f = parseEnvWithCounterAndIncrement(config.Env, c, store)
		}

		p := template{
			config.ID,
			config.PromptString,
			func(NextID string) func(string) string { return func(string) string { return NextID } }(config.NextID),
			func(ID string) func(error) string { return func(error) string { return ID } }(config.ID),
			f,
		}

		results = append(results, &p)
	}

	return results
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
