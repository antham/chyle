package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// decorater extends data from commit hashmap with data picked from third part service
type decorater interface {
	decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// decorate process all defined decorator and apply them
func decorate(decorators *map[string][]decorater, changelog *Changelog) (*Changelog, error) {
	var err error

	datas := []map[string]interface{}{}

	for _, d := range changelog.Datas {
		result := &d

		for _, decorator := range (*decorators)["datas"] {
			result, err = decorator.decorate(&d)

			if err != nil {
				return nil, err
			}
		}

		datas = append(datas, *result)
	}

	changelog.Datas = datas

	metadatas := changelog.Metadatas

	for _, decorator := range (*decorators)["metadatas"] {
		m, err := decorator.decorate(&metadatas)

		if err != nil {
			return nil, err
		}

		metadatas = *m
	}

	changelog.Metadatas = metadatas

	return changelog, nil
}

// createDatasDecorators build decorators dealing with datas extracted from repository
func createDatasDecorators(config *envh.EnvTree) (*[]decorater, error) {
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
			return nil, err
		}

		results = append(results, ex)
	}

	return &results, nil
}

// createDecorators build decorators from a config
func createDecorators(config *envh.EnvTree) (*map[string][]decorater, error) {
	datas, err := createDatasDecorators(config)

	if err != nil {
		return nil, err
	}

	return &map[string][]decorater{"metadatas": []decorater{}, "datas": *datas}, nil
}
