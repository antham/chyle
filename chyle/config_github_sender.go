package chyle

import (
	"github.com/antham/envh"
)

// githubSenderConfigurator validates github sender config defined through environment variables
type githubSenderConfigurator struct {
	config *envh.EnvTree
}

func (g *githubSenderConfigurator) process(config *CHYLE) (bool, error) {
	if g.isDisabled() {
		return false, nil
	}

	config.FEATURES.HASSENDERS = true
	config.FEATURES.HASGITHUBRELEASESENDER = true

	for _, f := range []func() error{
		g.validateCredentials,
		g.validateReleaseMandatoryFields,
		g.validateRepositoryName,
	} {
		if err := f(); err != nil {
			return false, err
		}
	}

	return false, nil
}

// isDisabled checks if github sender is enabled
func (g *githubSenderConfigurator) isDisabled() bool {
	return !g.config.IsExistingSubTree("CHYLE", "SENDERS", "GITHUB")
}

// validateCredentials checks github credentials to access remote api
func (g *githubSenderConfigurator) validateCredentials() error {
	return validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUB", "CREDENTIALS", "OAUTHTOKEN"}, {"CHYLE", "SENDERS", "GITHUB", "CREDENTIALS", "OWNER"}})
}

// validateReleaseMandatoryFields checks release mandatory field definition
func (g *githubSenderConfigurator) validateReleaseMandatoryFields() error {
	if err := validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUB", "RELEASE", "TAGNAME"}, {"CHYLE", "SENDERS", "GITHUB", "RELEASE", "TEMPLATE"}}); err != nil {
		return err
	}

	if err := validateTemplate(g.config, []string{"CHYLE", "SENDERS", "GITHUB", "RELEASE", "TEMPLATE"}); err != nil {
		return err
	}

	return nil
}

// validateRepositoryName checks if github repository name is defined
func (g *githubSenderConfigurator) validateRepositoryName() error {
	return validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUB", "REPOSITORY", "NAME"}})
}
