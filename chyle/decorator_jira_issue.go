package chyle

import (
	"net/http"
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

	return decorateMapFromJSONResponse(&j.client, req, chyleConfig.DECORATORS.JIRA.KEYS, commitMap)
}

// buildJiraIssueDecorator create a new jira ticket decorator
func buildJiraIssueDecorator() decorater {
	return jiraIssueDecorator{http.Client{}}
}
