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
	var ok bool

	if ID, ok = (*commitMap)["jiraIssueId"].(string); !ok {
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL+"/rest/api/2/issue/"+ID, nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME, chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD)
	req.Header.Set("Content-Type", "application/json")

	return jSONResponseDecorator{&j.client, req, chyleConfig.DECORATORS.JIRA.KEYS}.decorate(commitMap)
}

// buildJiraIssueDecorator create a new jira ticket decorator
func buildJiraIssueDecorator() decorater {
	return jiraIssueDecorator{http.Client{}}
}
