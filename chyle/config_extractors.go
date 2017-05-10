package chyle

import (
	"github.com/antham/envh"
)

// extractorsConfigurator validates jira config
// defined through environment variables
type extractorsConfigurator struct {
	config      *envh.EnvTree
	definedKeys []string
}

func (e *extractorsConfigurator) process(config *CHYLE) (bool, error) {
	if e.isDisabled() {
		return true, nil
	}

	config.FEATURES.HASEXTRACTORS = true

	for _, f := range []func() error{
		e.validateExtractors,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	e.setExtractors(config)

	return true, nil
}

// isDisabled checks if matchers are enabled
func (e *extractorsConfigurator) isDisabled() bool {
	return featureDisabled(e.config, [][]string{{"CHYLE", "EXTRACTORS"}})
}

// validateExtractors checks threesome extractor fields
func (e *extractorsConfigurator) validateExtractors() error {
	for _, key := range e.config.FindChildrenKeysUnsecured("CHYLE", "EXTRACTORS") {
		if err := validateEnvironmentVariablesDefinition(e.config, [][]string{{"CHYLE", "EXTRACTORS", key, "ORIGKEY"}, {"CHYLE", "EXTRACTORS", key, "DESTKEY"}, {"CHYLE", "EXTRACTORS", key, "REG"}}); err != nil {
			return err
		}

		if err := validateRegexp(e.config, []string{"CHYLE", "EXTRACTORS", key, "REG"}); err != nil {
			return err
		}

		e.definedKeys = append(e.definedKeys, key)
	}

	return nil
}

// setExtractors update chyleConfig with extracted extractors
func (e *extractorsConfigurator) setExtractors(config *CHYLE) {
	config.EXTRACTORS = map[string]map[string]string{}

	for _, key := range e.definedKeys {
		config.EXTRACTORS[key] = map[string]string{}

		for _, field := range []string{"ORIGKEY", "DESTKEY", "REG"} {
			config.EXTRACTORS[key][field] = e.config.FindStringUnsecured("CHYLE", "EXTRACTORS", key, field)
		}
	}
}
