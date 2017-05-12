package chyle

import (
	"github.com/antham/envh"
)

// githubReleaseSenderConfigurator validates github sender config defined through environment variables
type githubReleaseSenderConfigurator struct {
	config *envh.EnvTree
}

func (g *githubReleaseSenderConfigurator) process(config *CHYLE) (bool, error) {
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
func (g *githubReleaseSenderConfigurator) isDisabled() bool {
	return !g.config.IsExistingSubTree("CHYLE", "SENDERS", "GITHUBRELEASE")
}

// validateCredentials checks github credentials to access remote api
func (g *githubReleaseSenderConfigurator) validateCredentials() error {
	return validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUBRELEASE", "CREDENTIALS", "OAUTHTOKEN"}, {"CHYLE", "SENDERS", "GITHUBRELEASE", "CREDENTIALS", "OWNER"}})
}

// validateReleaseMandatoryFields checks release mandatory field definition
func (g *githubReleaseSenderConfigurator) validateReleaseMandatoryFields() error {
	if err := validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUBRELEASE", "RELEASE", "TAGNAME"}, {"CHYLE", "SENDERS", "GITHUBRELEASE", "RELEASE", "TEMPLATE"}}); err != nil {
		return err
	}

	if err := validateTemplate(g.config, []string{"CHYLE", "SENDERS", "GITHUBRELEASE", "RELEASE", "TEMPLATE"}); err != nil {
		return err
	}

	return nil
}

// validateRepositoryName checks if github repository name is defined
func (g *githubReleaseSenderConfigurator) validateRepositoryName() error {
	return validateEnvironmentVariablesDefinition(g.config, [][]string{{"CHYLE", "SENDERS", "GITHUBRELEASE", "REPOSITORY", "NAME"}})
}
