package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// codebeat:disable[TOO_MANY_IVARS]

// apiDecoratorConfig declares datas needed
// to validate an api configuration
type apiDecoratorConfig struct {
	extractorKey    string
	decoratorKey    string
	keysRef         *map[string]string
	credentialsRefs []struct {
		ref      *string
		keyChain []string
	}
	featureRefs []*bool
}

// codebeat:enable[TOO_MANY_IVARS]

// apiDecoratorConfigurator is a generic api
// decorator configurator it must be used with
// apiDecoratorConfig
type apiDecoratorConfigurator struct {
	config *envh.EnvTree
	apiDecoratorConfig
	definedKeys []string
}

func (a *apiDecoratorConfigurator) process(config *CHYLE) (bool, error) {
	if a.isDisabled() {
		return true, nil
	}

	for _, featureRef := range a.featureRefs {
		*featureRef = true
	}

	for _, f := range []func() error{
		a.validateCredentials,
		a.validateKeys,
		a.validateExtractor,
	} {
		if err := f(); err != nil {
			return true, err
		}
	}

	a.setKeys(config)
	a.setCredentials(config)

	return true, nil
}

// isDisabled checks if decorator is enabled
func (a *apiDecoratorConfigurator) isDisabled() bool {
	return featureDisabled(a.config, [][]string{
		{"CHYLE", "DECORATORS", a.decoratorKey},
		{"CHYLE", "EXTRACTORS", a.extractorKey},
	})
}

// validateExtractor checks if an extractor is defined to get
// data needed to contact remote api
func (a *apiDecoratorConfigurator) validateExtractor() error {
	return validateEnvironmentVariablesDefinition(a.config, [][]string{{"CHYLE", "EXTRACTORS", a.extractorKey, "ORIGKEY"}, {"CHYLE", "EXTRACTORS", a.extractorKey, "DESTKEY"}, {"CHYLE", "EXTRACTORS", a.extractorKey, "REG"}})
}

// validateCredentials checks credentials are defined to contact api
func (a *apiDecoratorConfigurator) validateCredentials() error {
	keyChains := [][]string{}

	for _, ref := range a.credentialsRefs {
		keyChains = append(keyChains, ref.keyChain)
	}

	if err := validateEnvironmentVariablesDefinition(a.config, keyChains); err != nil {
		return err
	}

	for _, keyChain := range keyChains {
		if keyChain[len(keyChain)-1] != "URL" {
			continue
		}

		if err := validateURL(a.config, keyChain); err != nil {
			return err
		}
	}

	return nil
}

// validateKeys checks key mapping between fields extracted from api and fields added to final struct
func (a *apiDecoratorConfigurator) validateKeys() error {
	keys, err := a.config.FindChildrenKeys("CHYLE", "DECORATORS", a.decoratorKey, "KEYS")

	if err != nil {
		return fmt.Errorf(`define at least one environment variable couple "CHYLE_DECORATORS_%s_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_%s_KEYS_*_FIELD", replace "*" with your own naming`, a.decoratorKey, a.decoratorKey)
	}

	for _, key := range keys {
		if err := validateEnvironmentVariablesDefinition(a.config, [][]string{{"CHYLE", "DECORATORS", a.decoratorKey, "KEYS", key, "DESTKEY"}, {"CHYLE", "DECORATORS", a.decoratorKey, "KEYS", key, "FIELD"}}); err != nil {
			return err
		}

		a.definedKeys = append(a.definedKeys, key)
	}

	return nil
}

// setCredentials update api credentials
func (a *apiDecoratorConfigurator) setCredentials(config *CHYLE) {
	for _, c := range a.credentialsRefs {
		*(c.ref) = a.config.FindStringUnsecured(c.keyChain...)
	}
}

// setKeys update keys needed for extraction
func (a *apiDecoratorConfigurator) setKeys(config *CHYLE) {
	ref := a.keysRef
	*ref = map[string]string{}

	for _, key := range a.definedKeys {
		datas := map[string]string{}

		for _, field := range []string{"DESTKEY", "FIELD"} {
			datas[field] = a.config.FindStringUnsecured(append([]string{"CHYLE", "DECORATORS", a.decoratorKey, "KEYS"}, key, field)...)
		}

		(*ref)[datas["DESTKEY"]] = datas["FIELD"]
	}
}
