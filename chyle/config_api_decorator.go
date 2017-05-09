package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

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
	return validateSubConfigPool(a.config, []string{"CHYLE", "EXTRACTORS", a.extractorKey}, []string{"ORIGKEY", "DESTKEY", "REG"})
}

// validateCredentials checks credentials are defined to contact api
func (a *apiDecoratorConfigurator) validateCredentials() error {
	if err := validateSubConfigPool(a.config, []string{"CHYLE", "DECORATORS", a.decoratorKey, "CREDENTIALS"}, []string{"URL", "USERNAME", "PASSWORD"}); err != nil {
		return err
	}

	if err := validateURL(a.config, []string{"CHYLE", "DECORATORS", a.decoratorKey, "CREDENTIALS", "URL"}); err != nil {
		return err
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
		if err := validateSubConfigPool(a.config, []string{"CHYLE", "DECORATORS", a.decoratorKey, "KEYS", key}, []string{"DESTKEY", "FIELD"}); err != nil {
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
