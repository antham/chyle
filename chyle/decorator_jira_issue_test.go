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

	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := newJiraIssueDecoratorFromPasswordAuth(*client)

	result, err := j.decorate(&map[string]interface{}{"test": "test", "jiraIssueId": "10000"})

	expected := map[string]interface{}{
		"test":         "test",
		"jiraIssueId":  "10000",
		"jiraIssueKey": "EX-1",
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestJiraDecoratorWithNoJiraIssueIdDefined(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := newJiraIssueDecoratorFromPasswordAuth(*client)

	result, err := j.decorate(&map[string]interface{}{"test": "test"})

	expected := map[string]interface{}{
		"test": "test",
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.False(t, gock.IsDone(), "Must have one pending request")
}
