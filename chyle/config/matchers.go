package config

import (
	"github.com/antham/chyle/chyle/matchers"

	"github.com/antham/envh"
)

// matchersConfigurator validates jira config
// defined through environment variables
type matchersConfigurator struct {
	config *envh.EnvTree
}

func (m *matchersConfigurator) process(config *CHYLE) (bool, error) {
	if m.isDisabled() {
		return true, nil
	}

	config.FEATURES.MATCHERS.ENABLED = true

	for _, f := range []func() error{
		m.validateRegexpMatchers,
		m.validateTypeMatcher,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	m.setMatchers(config)

	return true, nil
}

// isDisabled checks if matchers are enabled
func (m *matchersConfigurator) isDisabled() bool {
	return featureDisabled(m.config, [][]string{{"CHYLE", "MATCHERS"}})
}

// validateRegexpMatchers checks all config relying on valid regexp
func (m *matchersConfigurator) validateRegexpMatchers() error {
	for _, key := range []string{"MESSAGE", "COMMITTER", "AUTHOR"} {
		_, err := m.config.FindString("CHYLE", "MATCHERS", key)

		if err != nil {
			continue
		}

		if err := validateRegexp(m.config, []string{"CHYLE", "MATCHERS", key}); err != nil {
			return err
		}
	}

	return nil
}

// validateTypeMatcher checks custom field TYPE
func (m *matchersConfigurator) validateTypeMatcher() error {
	_, err := m.config.FindString("CHYLE", "MATCHERS", "TYPE")

	if err != nil {
		return nil
	}

	return validateOneOf(m.config, []string{"CHYLE", "MATCHERS", "TYPE"}, matchers.GetTypes())
}

// setMatchers update config with extracted matchers
func (m *matchersConfigurator) setMatchers(config *CHYLE) {
	config.MATCHERS = map[string]string{}

	for _, key := range m.config.FindChildrenKeysUnsecured("CHYLE", "MATCHERS") {
		config.MATCHERS[key] = m.config.FindStringUnsecured("CHYLE", "MATCHERS", key)
	}
}
