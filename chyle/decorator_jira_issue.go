package chyle

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/antham/envh"
)

// jiraIssueDecorator fetch data using jira issue api
type jiraIssueDecorator struct {
	client   http.Client
	username string
	password string
	URL      string
	keys     map[string]string
}

// newJiraIssueDecoratorFromPasswordAuth create a new jiraIssueDecorator
func newJiraIssueDecoratorFromPasswordAuth(client http.Client, username string, password string, URL string, keys map[string]string) (jiraIssueDecorator, error) {
	return jiraIssueDecorator{client, username, password, URL, keys}, nil
}

// decorate fetch remote jira service if a jiraIssueId is defined to fetch issue datas
func (j jiraIssueDecorator) decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID string

	if data, ok := (*commitMap)["jiraIssueId"]; true {
		if !ok {
			return commitMap, nil
		}

		if data, ok := data.(string); ok {
			ID = data
		}
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
		(*commitMap)[identifier] = nil

		if gjson.Get(buf.String(), key).Exists() {
			(*commitMap)[identifier] = gjson.Get(buf.String(), key).Value()
		}
	}

	return commitMap, nil
}

func buildJiraDecorator(config *envh.EnvTree) (decorater, error) {
	datas := map[string]string{}
	keyValues := map[string]string{}

	for _, k := range []string{"USERNAME", "PASSWORD", "URL"} {
		v, err := config.FindString("CREDENTIALS", k)

		if err != nil {
			return nil, fmt.Errorf(`"%s" variable not found in "JIRA" config`, k)
		}

		datas[k] = v
	}

	_, err := url.ParseRequestURI(datas["URL"])

	if err != nil {
		return nil, fmt.Errorf(`"%s" is not a valid absolute URL defined in "JIRA" config`, datas["URL"])
	}

	debug(`Decorator "USERNAME" defined with value "%s"`, datas["USERNAME"])
	debug(`Decorator "PASSWORD" defined`)
	debug(`Decorator "URL" defined with value "%s"`, datas["URL"])

	keys, err := config.FindChildrenKeys("KEYS")

	if err != nil {
		return nil, fmt.Errorf(`No "DECORATORS_JIRA_KEYS" key found`)
	}

	for _, k := range keys {
		key, err := config.FindString("KEYS", k, "DESTKEY")

		if err != nil {
			return nil, fmt.Errorf(`An environment variable suffixed with "DESTKEY" must be defined with "%s", like DECORATORS_JIRA_KEYS_%s_DESTKEY`, k, k)
		}

		value, err := config.FindString("KEYS", k, "FIELD")

		if err != nil {
			return nil, fmt.Errorf(`An environment variable suffixed with "FIELD" must be defined with "%s", like DECORATORS_JIRA_KEYS_%s_FIELD`, k, k)
		}

		debug(`Decorator KEY "%s" defined with value "%s"`, key, value)

		keyValues[key] = value
	}

	return newJiraIssueDecoratorFromPasswordAuth(http.Client{}, datas["USERNAME"], datas["PASSWORD"], datas["URL"], keyValues)
}
