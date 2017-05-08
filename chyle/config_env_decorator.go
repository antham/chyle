package chyle

import (
	"github.com/antham/envh"
)

// envDecoratorConfigurator validates environment variables decorator config
// defined through environment variables
type envDecoratorConfigurator struct {
	chyleConfig *CHYLE
	config      *envh.EnvTree
	definedKeys []string
}

func (e envDecoratorConfigurator) process() (bool, error) {
	if e.isDisabled() {
		return true, nil
	}

	for _, f := range []func() error{
		e.validateEnvironmentVariables,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	e.setEnvDecorator()

	return true, nil
}

// isDisabled checks if environment variable decorator is enabled
func (e envDecoratorConfigurator) isDisabled() bool {
	return featureDisabled(e.config, [][]string{{"CHYLE", "DECORATORS", "ENV"}})
}

// validateEnvironmentVariables checks env pairs are defined
func (e envDecoratorConfigurator) validateEnvironmentVariables() error {
	for _, key := range e.config.FindChildrenKeysUnsecured("CHYLE", "DECORATORS", "ENV") {
		if err := validateSubConfigPool(e.config, []string{"CHYLE", "DECORATORS", "ENV", key}, []string{"DESTKEY", "VARNAME"}); err != nil {
			return err
		}

		e.definedKeys = append(e.definedKeys, key)
	}

	return nil
}

// setEnvDecorator update decorator environment variables
func (e envDecoratorConfigurator) setEnvDecorator() {
	e.chyleConfig.DECORATORS.ENV = map[string]map[string]string{}

	for _, key := range e.definedKeys {
		e.chyleConfig.DECORATORS.ENV[key] = map[string]string{}

		for _, field := range []string{"DESTKEY", "VARNAME"} {
			e.chyleConfig.DECORATORS.ENV[key][field] = e.config.FindStringUnsecured("CHYLE", "DECORATORS", "ENV", key, field)
		}
	}
}
