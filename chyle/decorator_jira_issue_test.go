package chyle

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestJiraDecorator(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.DECORATORS.JIRA.KEYS = map[string]string{}
	chyleConfig.FEATURES.HASJIRADECORATOR = true
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL = "http://test.com"
	chyleConfig.DECORATORS.JIRA.KEYS["jiraIssueKey"] = "key"
	chyleConfig.DECORATORS.JIRA.KEYS["whatever"] = "whatever"

	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		JSON(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	gock.New("http://test.com/rest/api/2/issue/EX-1").
		Reply(200).
		JSON(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := jiraIssueDecorator{*client}

	// request with issue id
	result, err := j.decorate(&map[string]interface{}{"test": "test", "jiraIssueId": int64(10000)})

	expected := map[string]interface{}{
		"test":         "test",
		"jiraIssueId":  int64(10000),
		"jiraIssueKey": "EX-1",
		"whatever":     nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)

	// request with issue key
	result, err = j.decorate(&map[string]interface{}{"test": "test", "jiraIssueId": "EX-1"})

	expected = map[string]interface{}{
		"test":         "test",
		"jiraIssueId":  "EX-1",
		"jiraIssueKey": "EX-1",
		"whatever":     nil,
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestJiraDecoratorWithNoJiraIssueIdDefined(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		JSON(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := jiraIssueDecorator{*client}

	result, err := j.decorate(&map[string]interface{}{"test": "test"})

	expected := map[string]interface{}{
		"test": "test",
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.False(t, gock.IsDone(), "Must have one pending request")
}

func TestJiraDecoratorWhenIssueIsNotFound(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.DECORATORS.JIRA.KEYS = map[string]string{}
	chyleConfig.FEATURES.HASJIRADECORATOR = true
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL = "http://test.com"
	chyleConfig.DECORATORS.JIRA.KEYS["jiraIssueKey"] = "key"
	chyleConfig.DECORATORS.JIRA.KEYS["whatever"] = "whatever"

	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(404).
		JSON(`{"errorMessages":["Issue does not exist or you do not have permission to see it."],"errors":{}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := jiraIssueDecorator{*client}

	result, err := j.decorate(&map[string]interface{}{"test": "test", "jiraIssueId": int64(10000)})

	expected := map[string]interface{}{
		"test":        "test",
		"jiraIssueId": int64(10000),
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *result)
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}
