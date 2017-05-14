package config

import (
	"regexp"

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
	c := map[string]struct {
		re      **regexp.Regexp
		feature *bool
	}{
		"MESSAGE": {
			&config.MATCHERS.MESSAGE,
			&config.FEATURES.MATCHERS.MESSAGE,
		},
		"COMMITTER": {
			&config.MATCHERS.COMMITTER,
			&config.FEATURES.MATCHERS.COMMITTER,
		},
		"AUTHOR": {
			&config.MATCHERS.AUTHOR,
			&config.FEATURES.MATCHERS.AUTHOR,
		},
	}

	for _, key := range m.config.FindChildrenKeysUnsecured("CHYLE", "MATCHERS") {
		val := m.config.FindStringUnsecured("CHYLE", "MATCHERS", key)

		if key == "TYPE" {
			config.MATCHERS.TYPE = val
			config.FEATURES.MATCHERS.TYPE = true
		} else {
			*(c[key].re) = regexp.MustCompile(val)
			*(c[key].feature) = true
		}
	}
}
