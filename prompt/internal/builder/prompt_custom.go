package builder

import (
	"github.com/antham/strumt"
)

// NewPromptWithCustomHandlers creates a new prompts with the ability to customize onSuccess callback, onError  callback and parser
func NewPromptWithCustomHandlers(config EnvConfig, onSuccess func(string) string, onError func(error) string, parse func(string) error, store *Store) strumt.Prompter {
	return &template{
		config.ID,
		config.PromptString,
		onSuccess,
		onError,
		parse,
	}
}
