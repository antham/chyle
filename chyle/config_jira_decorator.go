package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// jiraDecoratorValidator validates jira config defined through environment variables
type jiraDecoratorValidator struct {
	config *envh.EnvTree
}

func (j jiraDecoratorValidator) validate() (bool, error) {
	if j.isDisabled() {
		return false, nil
	}

	for _, f := range []func() error{
		j.validateCredentials,
		j.validateKeys,
		j.validateExtractor,
	} {
		if err := f(); err != nil {
			return false, err
		}
	}

	return false, nil
}

// isDisabled checks if jira decorator is enabled
func (j jiraDecoratorValidator) isDisabled() bool {
	return featureDisabled(j.config, [][]string{
		{"CHYLE", "DECORATORS", "JIRA"},
		{"CHYLE", "EXTRACTORS", "JIRAISSUEID"},
	})
}

// validateExtractor checks if jira issue id extractor is defined
func (j jiraDecoratorValidator) validateExtractor() error {
	return validateSubConfigPool(j.config, []string{"CHYLE", "EXTRACTORS", "JIRAISSUEID"}, []string{"ORIGKEY", "DESTKEY", "REG"})
}

// validateCredentials checks jira credentials to access remote api
func (j jiraDecoratorValidator) validateCredentials() error {
	if err := validateSubConfigPool(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS"}, []string{"URL", "USERNAME", "PASSWORD"}); err != nil {
		return err
	}

	if err := validateURL(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "URL"}); err != nil {
		return err
	}

	return nil
}

// validateKeys checks key mapping between fields extracted from jira api and fields added to final struct
func (j jiraDecoratorValidator) validateKeys() error {
	keys, err := j.config.FindChildrenKeys("CHYLE", "DECORATORS", "JIRA", "KEYS")

	if err != nil {
		return fmt.Errorf(`define at least one environment variable couple "CHYLE_DECORATORS_JIRA_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_JIRA_KEYS_*_FIELD", replace "*" with your own naming`)
	}

	for _, key := range keys {
		if err := validateSubConfigPool(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "KEYS", key}, []string{"DESTKEY", "FIELD"}); err != nil {
			return err
		}
	}

	return nil
}
