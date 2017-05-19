package config

import (
	"github.com/antham/envh"
)

// customAPIDecoratorConfigurator creates a custom api configurater from apiDecoratorConfigurator
func customAPIDecoratorConfigurator(config *envh.EnvTree) configurater {
	return &apiDecoratorConfigurator{
		config: config,
		apiDecoratorConfig: apiDecoratorConfig{
			"CUSTOMAPIID",
			"CUSTOMAPI",
			&chyleConfig.DECORATORS.CUSTOMAPI.KEYS,
			[]struct {
				ref      *string
				keyChain []string
			}{
				{
					&chyleConfig.DECORATORS.CUSTOMAPI.CREDENTIALS.TOKEN,
					[]string{"CHYLE", "DECORATORS", "CUSTOMAPI", "CREDENTIALS", "TOKEN"},
				},
			},
			[]*bool{
				&chyleConfig.FEATURES.DECORATORS.ENABLED,
				&chyleConfig.FEATURES.DECORATORS.CUSTOMAPI,
			},
			[]func() error{},
			[]func(*CHYLE){},
		},
	}
}
