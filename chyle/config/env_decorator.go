package config

import (
	"github.com/antham/envh"
)

// envDecoratorConfigurator validates environment variables decorator config
// defined through environment variables
type envDecoratorConfigurator struct {
	config *envh.EnvTree
}

func (e *envDecoratorConfigurator) process(config *CHYLE) (bool, error) {
	if e.isDisabled() {
		return true, nil
	}

	config.FEATURES.DECORATORS.ENABLED = true
	config.FEATURES.DECORATORS.ENV = true

	for _, f := range []func() error{
		e.validateEnvironmentVariables,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	e.setEnvDecorator(config)

	return true, nil
}

// isDisabled checks if environment variable decorator is enabled
func (e *envDecoratorConfigurator) isDisabled() bool {
	return featureDisabled(e.config, [][]string{{"CHYLE", "DECORATORS", "ENV"}})
}

// validateEnvironmentVariables checks env pairs are defined
func (e *envDecoratorConfigurator) validateEnvironmentVariables() error {
	for _, key := range e.config.FindChildrenKeysUnsecured("CHYLE", "DECORATORS", "ENV") {
		if err := validateEnvironmentVariablesDefinition(e.config, [][]string{{"CHYLE", "DECORATORS", "ENV", key, "DESTKEY"}, {"CHYLE", "DECORATORS", "ENV", key, "VARNAME"}}); err != nil {
			return err
		}
	}

	return nil
}

// setEnvDecorator update decorator environment variables
func (e *envDecoratorConfigurator) setEnvDecorator(config *CHYLE) {
	config.DECORATORS.ENV = map[string]struct {
		DESTKEY string
		VARNAME string
	}{}

	for _, key := range e.config.FindChildrenKeysUnsecured("CHYLE", "DECORATORS", "ENV") {
		config.DECORATORS.ENV[key] = struct {
			DESTKEY string
			VARNAME string
		}{
			e.config.FindStringUnsecured("CHYLE", "DECORATORS", "ENV", key, "DESTKEY"),
			e.config.FindStringUnsecured("CHYLE", "DECORATORS", "ENV", key, "VARNAME"),
		}
	}
}
