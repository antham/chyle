package chyle

import (
	"bytes"
	"net/http"

	"github.com/tidwall/gjson"
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

// createDecorators build decorators from a config
func createDecorators() *map[string][]decorater {
	results := map[string][]decorater{"metadatas": {}, "datas": {}}

	if chyleConfig.FEATURES.HASJIRADECORATOR {
		results["datas"] = append(results["datas"], buildJiraIssueDecorator())
	}

	if chyleConfig.FEATURES.HASGITHUBDECORATOR {
		results["datas"] = append(results["datas"], buildGithubIssueDecorator())
	}

	if chyleConfig.FEATURES.HASENVDECORATOR {
		results["metadatas"] = append(results["metadatas"], buildEnvDecorators()...)
	}

	return &results
}

// jSONResponseDecorator extracts datas from a JSON api using defined keys
// and add it to final commitMap data structure
type jSONResponseDecorator struct {
	client  *http.Client
	request *http.Request
	pairs   map[string]string
}

// decorate fetch JSON datas and add the result to original commitMap array
func (j jSONResponseDecorator) decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	_, body, err := sendRequest(j.client, j.request)

	if err != nil {
		return commitMap, err
	}

	buf := bytes.NewBuffer(body)

	for identifier, key := range j.pairs {
		(*commitMap)[identifier] = nil

		if gjson.Get(buf.String(), key).Exists() {
			(*commitMap)[identifier] = gjson.Get(buf.String(), key).Value()
		}
	}

	return commitMap, nil
}
