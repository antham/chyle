package decorators

import (
	"fmt"
	"net/http"

	"github.com/antham/chyle/chyle/apih"
)

type githubIssueConfig struct {
	CREDENTIALS struct {
		OAUTHTOKEN string
		OWNER      string
	}
	REPOSITORY struct {
		NAME string
	}
	KEYS map[string]struct {
		DESTKEY string
		FIELD   string
	}
}

// githubIssueDecorator fetch data using github issue api
type githubIssueDecorator struct {
	client http.Client
	config githubIssueConfig
}

// decorate fetch remote github service if a github tikcet id is defined to fetch issue datas
func (g githubIssueDecorator) Decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID int64
	var ok bool

	if ID, ok = (*commitMap)["githubIssueId"].(int64); !ok {
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", g.config.CREDENTIALS.OWNER, g.config.REPOSITORY.NAME, ID), nil)

	if err != nil {
		return commitMap, err
	}

	apih.SetHeaders(req, map[string]string{
		"Authorization": "token " + g.config.CREDENTIALS.OAUTHTOKEN,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	return jSONResponseDecorator{&g.client, req, g.config.KEYS}.Decorate(commitMap)
}

// buildGithubIssueDecorator create a new github issue decorator
func buildGithubIssueDecorator(config githubIssueConfig) Decorater {
	return githubIssueDecorator{http.Client{}, config}
}
