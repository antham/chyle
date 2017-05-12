package chyle

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestGithubIssueDecorator(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.DECORATORS.GITHUB.KEYS = map[string]string{}
	chyleConfig.FEATURES.HASJIRADECORATOR = true
	chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
	chyleConfig.DECORATORS.GITHUB.CREDENTIALS.OWNER = "user"
	chyleConfig.DECORATORS.GITHUB.REPOSITORY.NAME = "repository"
	chyleConfig.DECORATORS.GITHUB.KEYS["milestoneCreator"] = "milestone.creator.id"
	chyleConfig.DECORATORS.GITHUB.KEYS["whatever"] = "whatever"

	defer gock.Off()

	issueResponse, err := ioutil.ReadFile("fixtures/3-github-issue-fetch-response.json")

	assert.NoError(t, err, "Must read json fixture file")

	gock.New("https://api.github.com/repos/user/repository/issues/10000").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		Reply(200).
		JSON(string(issueResponse))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := githubIssueDecorator{*client}

	result, err := j.decorate(&map[string]interface{}{"test": "test", "githubIssueId": int64(10000)})

	expected := map[string]interface{}{
		"test":             "test",
		"githubIssueId":    int64(10000),
		"milestoneCreator": float64(1),
		"whatever":         nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubDecoratorWithNoGithubIssueIdDefined(t *testing.T) {
	defer gock.Off()

	issueResponse, err := ioutil.ReadFile("fixtures/3-github-issue-fetch-response.json")

	assert.NoError(t, err, "Must read json fixture file")

	gock.New("https://api.github.com/repos/user/repository/issues/10000").
		Reply(200).
		JSON(string(issueResponse))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := githubIssueDecorator{*client}

	result, err := j.decorate(&map[string]interface{}{"test": "test"})

	expected := map[string]interface{}{
		"test": "test",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)
	assert.False(t, gock.IsDone(), "Must have one pending request")
}
