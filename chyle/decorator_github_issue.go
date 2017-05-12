package chyle

import (
	"fmt"
	"net/http"
)

// githubIssueDecorator fetch data using github issue api
type githubIssueDecorator struct {
	client http.Client
}

// decorate fetch remote github service if a github tikcet id is defined to fetch issue datas
func (g githubIssueDecorator) decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID int64
	var ok bool

	if ID, ok = (*commitMap)["githubIssueId"].(int64); !ok {
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OWNER, chyleConfig.DECORATORS.GITHUB.REPOSITORY.NAME, ID), nil)

	if err != nil {
		return commitMap, err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OAUTHTOKEN,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	return jSONResponseDecorator{&g.client, req, chyleConfig.DECORATORS.GITHUB.KEYS}.decorate(commitMap)
}

// buildGithubIssueDecorator create a new github issue decorator
func buildGithubIssueDecorator() decorater {
	return githubIssueDecorator{http.Client{}}
}
