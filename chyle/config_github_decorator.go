package chyle

import (
	"github.com/antham/envh"
)

// githubDecoratorConfigurator creates a github configurater from apiDecoratorConfigurator
func githubDecoratorConfigurator(config *envh.EnvTree) configurater {
	return &apiDecoratorConfigurator{
		config: config,
		apiDecoratorConfig: apiDecoratorConfig{
			"GITHUBISSUEID",
			"GITHUB",
			&chyleConfig.DECORATORS.GITHUB.KEYS,
			[]struct {
				ref      *string
				keyChain []string
			}{
				{
					&chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OAUTHTOKEN,
					[]string{"CHYLE", "DECORATORS", "GITHUB", "CREDENTIALS", "OAUTHTOKEN"},
				},
				{
					&chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OWNER,
					[]string{"CHYLE", "DECORATORS", "GITHUB", "CREDENTIALS", "OWNER"},
				},
				{
					&chyleConfig.DECORATORS.GITHUB.REPOSITORY.NAME,
					[]string{"CHYLE", "DECORATORS", "GITHUB", "REPOSITORY", "NAME"},
				},
			},
			[]*bool{
				&chyleConfig.FEATURES.HASDECORATORS,
				&chyleConfig.FEATURES.HASGITHUBDECORATOR,
			},
		},
	}
}
