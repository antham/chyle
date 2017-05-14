package decorators

import (
	"bytes"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/antham/chyle/chyle/apih"
	"github.com/antham/chyle/chyle/types"
)

// Decorater extends data from commit hashmap with data picked from third part service
type Decorater interface {
	Decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// Decorate process all defined decorator and apply them
func Decorate(decorators *map[string][]Decorater, changelog *types.Changelog) (*types.Changelog, error) {
	var err error

	datas := []map[string]interface{}{}

	for _, d := range changelog.Datas {
		result := &d

		for _, decorator := range (*decorators)["datas"] {
			result, err = decorator.Decorate(&d)

			if err != nil {
				return nil, err
			}
		}

		datas = append(datas, *result)
	}

	changelog.Datas = datas

	metadatas := changelog.Metadatas

	for _, decorator := range (*decorators)["metadatas"] {
		m, err := decorator.Decorate(&metadatas)

		if err != nil {
			return nil, err
		}

		metadatas = *m
	}

	changelog.Metadatas = metadatas

	return changelog, nil
}

// CreateDecorators build decorators from a config
func CreateDecorators(features Features, decorators Config) *map[string][]Decorater {
	results := map[string][]Decorater{"metadatas": {}, "datas": {}}

	if features.JIRAISSUE {
		results["datas"] = append(results["datas"], buildJiraIssueDecorator(decorators.JIRAISSUE))
	}

	if features.GITHUBISSUE {
		results["datas"] = append(results["datas"], buildGithubIssueDecorator(decorators.GITHUBISSUE))
	}

	if features.ENV {
		results["metadatas"] = append(results["metadatas"], buildEnvDecorators(decorators.ENV)...)
	}

	return &results
}

// jSONResponseDecorator extracts datas from a JSON api using defined keys
// and add it to final commitMap data structure
type jSONResponseDecorator struct {
	client  *http.Client
	request *http.Request
	pairs   map[string]struct {
		DESTKEY string
		FIELD   string
	}
}

// decorate fetch JSON datas and add the result to original commitMap array
func (j jSONResponseDecorator) Decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	statusCode, body, err := apih.SendRequest(j.client, j.request)

	if statusCode == 404 {
		return commitMap, nil
	}

	if err != nil {
		return commitMap, err
	}

	buf := bytes.NewBuffer(body)

	for _, pair := range j.pairs {
		(*commitMap)[pair.DESTKEY] = nil

		if gjson.Get(buf.String(), pair.FIELD).Exists() {
			(*commitMap)[pair.DESTKEY] = gjson.Get(buf.String(), pair.FIELD).Value()
		}
	}

	return commitMap, nil
}
