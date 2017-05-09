package chyle

import (
	"github.com/antham/envh"
)

// jiraDecoratorConfigurator creates a jira configurater from apiDecoratorConfigurator
func jiraDecoratorConfigurator(config *envh.EnvTree) configurater {
	return &apiDecoratorConfigurator{
		config: config,
		apiDecoratorConfig: apiDecoratorConfig{
			"JIRAISSUEID",
			"JIRA",
			&chyleConfig.DECORATORS.JIRA.KEYS,
			[]struct {
				ref      *string
				keyChain []string
			}{
				{
					&chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL,
					[]string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "URL"},
				},
				{
					&chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME,
					[]string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "USERNAME"},
				},
				{
					&chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD,
					[]string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "PASSWORD"},
				},
			},
			[]*bool{
				&chyleConfig.FEATURES.HASDECORATORS,
				&chyleConfig.FEATURES.HASJIRADECORATOR,
			},
		},
	}
}
