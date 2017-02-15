package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// Expander extends data from commit hashmap with data picked from third part service
type Expander interface {
	Expand(*map[string]interface{}) (*map[string]interface{}, error)
}

// Expand process all defined expander and apply them against every commit map
func Expand(expanders *[]Expander, commitMaps *[]map[string]interface{}) (*[]map[string]interface{}, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, expander := range *expanders {
			result, err = expander.Expand(&commitMap)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &results, nil
}

// CreateExpanders build expanders from a config
func CreateExpanders(config *envh.EnvTree) (*[]Expander, error) {
	results := []Expander{}

	var ex Expander
	var err error
	var subConfig envh.EnvTree

	for _, k := range config.GetChildrenKeys() {
		switch k {
		case "JIRA":
			subConfig, err = config.FindSubTree("JIRA")

			if err != nil {
				break
			}

			ex, err = buildJiraExpander(&subConfig)
		default:
			err = fmt.Errorf(`a wrong expander key containing "%s" was defined`, k)
		}

		if err != nil {
			return &[]Expander{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
