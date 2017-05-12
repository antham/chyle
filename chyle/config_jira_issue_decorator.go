package chyle

import (
	"github.com/antham/envh"
)

// jiraIssueDecoratorConfigurator creates a jira configurater from apiDecoratorConfigurator
func jiraIssueDecoratorConfigurator(config *envh.EnvTree) configurater {
	return &apiDecoratorConfigurator{
		config: config,
		apiDecoratorConfig: apiDecoratorConfig{
			"JIRAISSUEID",
			"JIRAISSUE",
			&chyleConfig.DECORATORS.JIRAISSUE.KEYS,
			[]struct {
				ref      *string
				keyChain []string
			}{
				{
					&chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.URL,
					[]string{"CHYLE", "DECORATORS", "JIRAISSUE", "CREDENTIALS", "URL"},
				},
				{
					&chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.USERNAME,
					[]string{"CHYLE", "DECORATORS", "JIRAISSUE", "CREDENTIALS", "USERNAME"},
				},
				{
					&chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.PASSWORD,
					[]string{"CHYLE", "DECORATORS", "JIRAISSUE", "CREDENTIALS", "PASSWORD"},
				},
			},
			[]*bool{
				&chyleConfig.FEATURES.HASDECORATORS,
				&chyleConfig.FEATURES.HASJIRAISSUEDECORATOR,
			},
		},
	}
}
