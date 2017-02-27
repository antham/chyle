package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// decorater extends data from commit hashmap with data picked from third part service
type decorater interface {
	decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// decorate process all defined decorator and apply them against every commit map
func decorate(decorators *[]decorater, commitMaps *[]map[string]interface{}) (*[]map[string]interface{}, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, decorator := range *decorators {
			result, err = decorator.decorate(&commitMap)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &results, nil
}

// createDecorators build decorators from a config
func createDecorators(config *envh.EnvTree) (*[]decorater, error) {
	results := []decorater{}

	var ex decorater
	var err error
	var subConfig envh.EnvTree

	for _, k := range config.GetChildrenKeys() {
		switch k {
		case "JIRA":
			subConfig, err = config.FindSubTree("JIRA")

			if err != nil {
				break
			}

			ex, err = buildJiraDecorator(&subConfig)
		default:
			err = fmt.Errorf(`a wrong decorator key containing "%s" was defined`, k)
		}

		if err != nil {
			return &[]decorater{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
