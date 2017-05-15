package decorators

import (
	"os"
)

type envConfig map[string]struct {
	DESTKEY string
	VARNAME string
}

// env dumps an environment variable into metadatas
type env struct {
	varName string
	destKey string
}

// Decorate adds an environment variable to changelog metadatas
func (e env) Decorate(metadatas *map[string]interface{}) (*map[string]interface{}, error) {
	(*metadatas)[e.destKey] = os.Getenv(e.varName)

	return metadatas, nil
}

// buildEnvs creates a list of env decorators
func buildEnvs(envs envConfig) []Decorater {
	results := []Decorater{}

	for _, e := range envs {
		results = append(results, env{
			e.VARNAME,
			e.DESTKEY,
		})
	}

	return results
}
