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
				{"1", "Add a type matcher", "matcherType"},
				{"2", "Add a message matcher", "matcherMessage"},
				{"3", "Add a committer matcher", "matcherCommitter"},
				{"4", "Add an author matcher", "matcherAuthor"},
			},
		),
	),
}

var matcher = []builder.EnvConfig{
	{"matcherType", "matcherChoice", "CHYLE_MATCHERS_TYPE", "Enter a matcher type (regular or merge)"},
	{"matcherMessage", "matcherChoice", "CHYLE_MATCHERS_MESSAGE", "Enter a regexp to match commit message"},
	{"matcherCommitter", "matcherChoice", "CHYLE_MATCHERS_COMMITTER", "Enter a regexp to match git committer"},
	{"matcherAuthor", "matcherChoice", "CHYLE_MATCHERS_AUTHOR", "Enter a regexp to match git author"},
}
