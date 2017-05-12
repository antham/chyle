package chyle

import (
	"fmt"
	"net/http"
)

// jiraIssueDecorator fetch data using jira issue api
type jiraIssueDecorator struct {
	client http.Client
}

// decorate fetch remote jira service if a jiraIssueId is defined to fetch issue datas
func (j jiraIssueDecorator) decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var URLpattern string

	switch (*commitMap)["jiraIssueId"].(type) {
	case string:
		URLpattern = "%s/rest/api/2/issue/%s"
	case int64:
		URLpattern = "%s/rest/api/2/issue/%d"
	default:
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(URLpattern, chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.URL, (*commitMap)["jiraIssueId"]), nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.USERNAME, chyleConfig.DECORATORS.JIRAISSUE.CREDENTIALS.PASSWORD)
	req.Header.Set("Content-Type", "application/json")

	return jSONResponseDecorator{&j.client, req, chyleConfig.DECORATORS.JIRAISSUE.KEYS}.decorate(commitMap)
}

// buildJiraIssueDecorator create a new jira ticket decorator
func buildJiraIssueDecorator() decorater {
	return jiraIssueDecorator{http.Client{}}
}
