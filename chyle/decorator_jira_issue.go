package chyle

import (
	"bytes"
	"net/http"

	"github.com/tidwall/gjson"
)

// jiraIssueDecorator fetch data using jira issue api
type jiraIssueDecorator struct {
	client http.Client
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

	req, err := http.NewRequest("GET", chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL+"/rest/api/2/issue/"+ID, nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME, chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD)
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

	for identifier, key := range chyleConfig.DECORATORS.JIRA.KEYS {
		(*commitMap)[identifier] = nil

		if gjson.Get(buf.String(), key).Exists() {
			(*commitMap)[identifier] = gjson.Get(buf.String(), key).Value()
		}
	}

	return commitMap, nil
}

// buildJiraDecorator create a new jira ticket decorator
func buildJiraDecorator() decorater {
	return jiraIssueDecorator{http.Client{}}
}
