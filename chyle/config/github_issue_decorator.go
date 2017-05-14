package config

import (
	"github.com/antham/envh"
)

// githubIssueDecoratorConfigurator creates a github configurater from apiDecoratorConfigurator
func githubIssueDecoratorConfigurator(config *envh.EnvTree) configurater {
	return &apiDecoratorConfigurator{
		config: config,
		apiDecoratorConfig: apiDecoratorConfig{
			"GITHUBISSUEID",
			"GITHUBISSUE",
			&chyleConfig.DECORATORS.GITHUBISSUE.KEYS,
			[]struct {
				ref      *string
				keyChain []string
			}{
				{
					&chyleConfig.DECORATORS.GITHUBISSUE.CREDENTIALS.OAUTHTOKEN,
					[]string{"CHYLE", "DECORATORS", "GITHUBISSUE", "CREDENTIALS", "OAUTHTOKEN"},
				},
				{
					&chyleConfig.DECORATORS.GITHUBISSUE.CREDENTIALS.OWNER,
					[]string{"CHYLE", "DECORATORS", "GITHUBISSUE", "CREDENTIALS", "OWNER"},
				},
				{
					&chyleConfig.DECORATORS.GITHUBISSUE.REPOSITORY.NAME,
					[]string{"CHYLE", "DECORATORS", "GITHUBISSUE", "REPOSITORY", "NAME"},
				},
			},
			[]*bool{
				&chyleConfig.FEATURES.DECORATORS.ENABLED,
				&chyleConfig.FEATURES.DECORATORS.GITHUBISSUE,
			},
		},
	}
}
