package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// jiraDecoratorProcessor validates jira config
// defined through environment variables
type jiraDecoratorProcessor struct {
	chyleConfig *CHYLE
	config      *envh.EnvTree
	definedKeys []string
}

func (j jiraDecoratorProcessor) process() (bool, error) {
	if j.isDisabled() {
		return true, nil
	}

	for _, f := range []func() error{
		j.validateCredentials,
		j.validateKeys,
		j.validateExtractor,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	j.setKeys()

	return true, nil
}

// isDisabled checks if jira decorator is enabled
func (j jiraDecoratorProcessor) isDisabled() bool {
	return featureDisabled(j.config, [][]string{
		{"CHYLE", "DECORATORS", "JIRA"},
		{"CHYLE", "EXTRACTORS", "JIRAISSUEID"},
	})
}

// validateExtractor checks if jira issue id extractor is defined
func (j jiraDecoratorProcessor) validateExtractor() error {
	return validateSubConfigPool(j.config, []string{"CHYLE", "EXTRACTORS", "JIRAISSUEID"}, []string{"ORIGKEY", "DESTKEY", "REG"})
}

// validateCredentials checks jira credentials to access remote api
func (j jiraDecoratorProcessor) validateCredentials() error {
	if err := validateSubConfigPool(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS"}, []string{"URL", "USERNAME", "PASSWORD"}); err != nil {
		return err
	}

	if err := validateURL(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "URL"}); err != nil {
		return err
	}

	return nil
}

// validateKeys checks key mapping between fields extracted from jira api and fields added to final struct
func (j jiraDecoratorProcessor) validateKeys() error {
	keys, err := j.config.FindChildrenKeys("CHYLE", "DECORATORS", "JIRA", "KEYS")

	if err != nil {
		return fmt.Errorf(`define at least one environment variable couple "CHYLE_DECORATORS_JIRA_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_JIRA_KEYS_*_FIELD", replace "*" with your own naming`)
	}

	for _, key := range keys {
		if err := validateSubConfigPool(j.config, []string{"CHYLE", "DECORATORS", "JIRA", "KEYS", key}, []string{"DESTKEY", "FIELD"}); err != nil {
			return err
		}

		j.definedKeys = append(j.definedKeys, key)
	}

	return nil
}

// setKeys update jira keys
func (j jiraDecoratorProcessor) setKeys() {
	j.chyleConfig.DECORATORS.JIRA.KEYS = map[string]string{}

	for _, key := range j.definedKeys {
		datas := map[string]string{}

		for _, field := range []string{"DESTKEY", "FIELD"} {
			datas[field] = j.config.FindStringUnsecured(append([]string{"CHYLE", "DECORATORS", "JIRA", "KEYS"}, key, field)...)
		}

		j.chyleConfig.DECORATORS.JIRA.KEYS[datas["DESTKEY"]] = datas["FIELD"]
	}
}
