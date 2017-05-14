package decorators

import (
	"os"
)

type envConfig map[string]struct {
	DESTKEY string
	VARNAME string
}

// envDecorator dump an environment variable into metadatas
type envDecorator struct {
	varName string
	destKey string
}

// Decorate add environment variable to changelog metadatas
func (e envDecorator) Decorate(metadatas *map[string]interface{}) (*map[string]interface{}, error) {
	(*metadatas)[e.destKey] = os.Getenv(e.varName)

	return metadatas, nil
}

// buildEnvDecorators creates a list of env decorators
func buildEnvDecorators(envs envConfig) []Decorater {
	results := []Decorater{}

	for _, e := range envs {
		results = append(results, envDecorator{
			e.VARNAME,
			e.DESTKEY,
		})
	}

	return results
}
