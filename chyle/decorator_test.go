package chyle

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestDecorator(t *testing.T) {
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

	gock.New("http://test.com/rest/api/2/issue/ABC-123").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10001","self":"http://test.com/jira/rest/api/2/issue/10001","key":"ABC-123","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := newJiraIssueDecoratorFromPasswordAuth(*client)

	decorators := map[string][]decorater{
		"datas":     {j},
		"metadatas": {},
	}

	changelog := Changelog{
		Datas: []map[string]interface{}{
			{
				"test":        "test1",
				"jiraIssueId": "10000",
			},
			{
				"test":        "test2",
				"jiraIssueId": "ABC-123",
			}},
		Metadatas: map[string]interface{}{},
	}

	result, err := decorate(&decorators, &changelog)

	expected := Changelog{
		Datas: []map[string]interface{}{
			{
				"test":         "test1",
				"jiraIssueId":  "10000",
				"jiraIssueKey": "EX-1",
			},
			{
				"test":         "test2",
				"jiraIssueId":  "ABC-123",
				"jiraIssueKey": "ABC-123",
			}},
		Metadatas: map[string]interface{}{},
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestCreateDecorators(t *testing.T) {
	chyleConfig.DECORATORS.JIRA.KEYS = map[string]string{}

	chyleConfig.FEATURES.HASJIRADECORATOR = true
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.USERNAME = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.PASSWORD = "test"
	chyleConfig.DECORATORS.JIRA.CREDENTIALS.URL = "http://test.com"
	chyleConfig.DECORATORS.JIRA.KEYS["jiraTicketDescription"] = "fields.summary"

	d := createDecorators()

	assert.Len(t, (*d)["datas"], 1, "Must return 1 decorator")
}
