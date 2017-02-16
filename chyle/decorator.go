package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// Decorater extends data from commit hashmap with data picked from third part service
type Decorater interface {
	Decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// Decorate process all defined decorator and apply them against every commit map
func Decorate(decorators *[]Decorater, commitMaps *[]map[string]interface{}) (*[]map[string]interface{}, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, decorator := range *decorators {
			result, err = decorator.Decorate(&commitMap)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &results, nil
}

// CreateDecorators build decorators from a config
func CreateDecorators(config *envh.EnvTree) (*[]Decorater, error) {
	results := []Decorater{}

	var ex Decorater
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
			return &[]Decorater{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
