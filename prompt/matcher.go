package prompt

import (
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
				builder.SwitchConfig{"1", "Add a type matcher", "matcherType"},
				builder.SwitchConfig{"2", "Add a message matcher", "matcherMessage"},
				builder.SwitchConfig{"3", "Add a committer matcher", "matcherCommitter"},
				builder.SwitchConfig{"4", "Add an author matcher", "matcherAuthor"},
			},
		),
	),
}

var matcher = []builder.EnvConfig{
	builder.EnvConfig{"matcherType", "matcherChoice", "CHYLE_MATCHERS_TYPE", "Enter a matcher type (regular or merge)"},
	builder.EnvConfig{"matcherMessage", "matcherChoice", "CHYLE_MATCHERS_MESSAGE", "Enter a regexp to match commit message"},
	builder.EnvConfig{"matcherCommitter", "matcherChoice", "CHYLE_MATCHERS_COMMITTER", "Enter a regexp to match git committer"},
	builder.EnvConfig{"matcherAuthor", "matcherChoice", "CHYLE_MATCHERS_AUTHOR", "Enter a regexp to match git author"},
}
