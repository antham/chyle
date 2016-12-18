package chyle

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// Expander extends data from commit hashmap with data picked from third part service
type Expander interface {
	Expand(*map[string]interface{}) (*map[string]interface{}, error)
}

// JiraIssueExpander fetch data using jira issue api
type JiraIssueExpander struct {
	client   http.Client
	username string
	password string
	URL      string
	keys     map[string]string
}

// NewJiraIssueExpanderFromPasswordAuth create a new JiraIssueExpander
func NewJiraIssueExpanderFromPasswordAuth(client http.Client, username string, password string, URL string, keys map[string]string) (JiraIssueExpander, error) {
	return JiraIssueExpander{client, username, password, URL, keys}, nil
}

// Expand fetch remote jira service if a jiraIssueId is defined to fetch issue datas
func (j JiraIssueExpander) Expand(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID string

	if data, ok := (*commitMap)["jiraIssueId"]; ok {
		if data, ok := data.(string); ok {
			ID = data
		}
	}

	if ID == "" {
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", j.URL+"/rest/api/2/issue/"+ID, nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(j.username, j.password)
	req.Header.Set("Content-Type", "application/json")

	rep, err := j.client.Do(req)

	if err != nil {
		return commitMap, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = rep.Write(buf)

	if err != nil {
		return commitMap, err
	}

	for identifier, key := range j.keys {
		if gjson.Get(buf.String(), key).Exists() {
			(*commitMap)[identifier] = gjson.Get(buf.String(), key).Value()
		} else {
			(*commitMap)[identifier] = nil
		}
	}

	return commitMap, nil
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

func buildJiraExpander(credentials map[string]string, keys map[string]string) (Expander, error) {
	var URL *url.URL

	for _, k := range []string{"username", "password", "url"} {
		if _, ok := credentials[k]; !ok {
			return nil, fmt.Errorf(`"%s" must be defined in jira config`, k)
		}
	}

	URL, err := url.Parse(credentials["url"])

	if err != nil {
		return nil, fmt.Errorf(`"%s" not a valid URL defined in jira config`, credentials["url"])
	}

	return NewJiraIssueExpanderFromPasswordAuth(http.Client{}, credentials["username"], credentials["password"], URL.String(), keys)
}

// CreateExpanders build expanders from a config
func CreateExpanders(config *viper.Viper) (*[]Expander, error) {
	results := []Expander{}

	for k := range config.GetStringMap("expanders") {
		var ex Expander
		var err error

		switch k {
		case "jira":
			if !config.IsSet("expanders.jira.credentials") || !config.IsSet("expanders.jira.keys") {
				return nil, fmt.Errorf(`"credentials" and "keys" key must be defined`)
			}

			ex, err = buildJiraExpander(config.GetStringMapString("expanders.jira.credentials"), config.GetStringMapString("expanders.jira.keys"))
		default:
			err = fmt.Errorf(`"%s" is not a valid expander structure`, k)
		}

		if err != nil {
			return &[]Expander{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
