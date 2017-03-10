package chyle

import (
	"fmt"
	"os"

	"github.com/antham/envh"
)

// envDecorator dump an environment variable into metadatas
type envDecorator struct {
	envVar  string
	destKey string
}

// decorate add environment variable to changelog metadatas
func (e envDecorator) decorate(metadatas *map[string]interface{}) (*map[string]interface{}, error) {
	(*metadatas)[os.Getenv(e.destKey)] = os.Getenv(e.envVar)

	return metadatas, nil
}

// buildEnvDecorators creates a list of env decorators
func buildEnvDecorators(config *envh.EnvTree) ([]decorater, error) {
	var err error

	results := []decorater{}

	for _, identifier := range config.GetChildrenKeys() {
		e := map[string]string{}

		for _, v := range []string{"VALUE", "DESTKEY"} {
			e[v], err = config.FindString(identifier, v)

			if err != nil {
				return nil, fmt.Errorf(`An environment variable suffixed with "%s" must be defined with "%s", like DECORATORS_ENV_%s_%s`, v, identifier, identifier, v)
			}
		}

		debug(`Decorator "%s" "VALUE" defined with value "%s"`, identifier, e["VALUE"])
		debug(`Decorator "%s" "DESTKEY" defined with value "%s"`, identifier, e["DESTKEY"])

		results = append(results, envDecorator{
			e["VALUE"],
			e["DESTKEY"],
		})
	}

	return results, nil
}
