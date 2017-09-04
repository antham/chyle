package prompt

import (
	"github.com/antham/strumt"
)

func newMatchers(store *Store) []strumt.Prompter {
	return mergePrompters(
		matcherChoice,
		newEnvPrompts(matcher, store),
	)
}

var matcherChoice = []strumt.Prompter{
	&switchPrompt{
		"matcherChoice",
		addMainMenuAndQuitChoice(
			[]switchChoice{
				switchChoice{"1", "Add a type matcher", "matcherType"},
				switchChoice{"2", "Add a message matcher", "matcherMessage"},
				switchChoice{"3", "Add a committer matcher", "matcherCommitter"},
				switchChoice{"4", "Add an author matcher", "matcherAuthor"},
			},
		),
	},
}

var matcher = []envConfig{
	envConfig{"matcherType", "matcherChoice", "CHYLE_MATCHERS_TYPE", "Enter a matcher type (regular or merge)"},
	envConfig{"matcherMessage", "matcherChoice", "CHYLE_MATCHERS_MESSAGE", "Enter a regexp to match commit message"},
	envConfig{"matcherCommitter", "matcherChoice", "CHYLE_MATCHERS_COMMITTER", "Enter a regexp to match git committer"},
	envConfig{"matcherAuthor", "matcherChoice", "CHYLE_MATCHERS_AUTHOR", "Enter a regexp to match git author"},
}
