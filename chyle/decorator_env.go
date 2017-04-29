package chyle

import (
	"os"
)

// envDecorator dump an environment variable into metadatas
type envDecorator struct {
	envVar  string
	destKey string
}

// decorate add environment variable to changelog metadatas
func (e envDecorator) decorate(metadatas *map[string]interface{}) (*map[string]interface{}, error) {
	(*metadatas)[e.destKey] = os.Getenv(e.envVar)

	return metadatas, nil
}

// buildEnvDecorators creates a list of env decorators
func buildEnvDecorators() []decorater {
	results := []decorater{}

	for _, e := range chyleConfig.DECORATORS.ENV {
		results = append(results, envDecorator{
			e["VALUE"],
			e["DESTKEY"],
		})
	}

	return results
}
