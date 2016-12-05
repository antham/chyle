package chyle

import (
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestJiraExpander(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	jiraClient, err := jira.NewClient(client, "http://test.com")

	if err != nil {
		logrus.Fatal(err)
	}

	j, err := NewJiraIssueExpanderFromPasswordAuth("test", "test", "http://test.com")

	assert.NoError(t, err, "Must return no errors")

	j.client = jiraClient

	result, err := j.Expand(&map[string]interface{}{"test": "test", "jiraIssueId": "10000"})

	expected := map[string]interface{}{
		"test":        "test",
		"jiraIssueId": "10000",
		"jiraIssue": &jira.Issue{
			Expand: "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations",
			ID:     "10000",
			Self:   "http://test.com/jira/rest/api/2/issue/10000",
			Key:    "EX-1",
		}}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestJiraExpanderWithNoJiraIssueIdDefined(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	jiraClient, err := jira.NewClient(client, "http://test.com")

	if err != nil {
		logrus.Fatal(err)
	}

	j, err := NewJiraIssueExpanderFromPasswordAuth("test", "test", "http://test.com")

	assert.NoError(t, err, "Must return no errors")

	j.client = jiraClient

	result, err := j.Expand(&map[string]interface{}{"test": "test"})

	expected := map[string]interface{}{
		"test": "test",
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.False(t, gock.IsDone(), "Must have one pending request")

}

func TestExpander(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	gock.New("http://test.com/rest/api/2/issue/ABC-123").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10001","self":"http://test.com/jira/rest/api/2/issue/10001","key":"ABC-123","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	jiraClient, err := jira.NewClient(client, "http://test.com")

	if err != nil {
		logrus.Fatal(err)
	}

	j, err := NewJiraIssueExpanderFromPasswordAuth("test", "test", "http://test.com")

	assert.NoError(t, err, "Must return no errors")

	j.client = jiraClient

	expanders := []Expander{
		j,
	}

	commitMaps := []map[string]interface{}{
		map[string]interface{}{
			"test":        "test1",
			"jiraIssueId": "10000",
		},
		map[string]interface{}{
			"test":        "test2",
			"jiraIssueId": "ABC-123",
		},
	}

	result, err := Expand(&expanders, &commitMaps)

	expected := []map[string]interface{}{
		map[string]interface{}{
			"test":        "test1",
			"jiraIssueId": "10000",
			"jiraIssue": &jira.Issue{
				Expand: "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations",
				ID:     "10000",
				Self:   "http://test.com/jira/rest/api/2/issue/10000",
				Key:    "EX-1",
			},
		},
		map[string]interface{}{
			"test":        "test2",
			"jiraIssueId": "ABC-123",
			"jiraIssue": &jira.Issue{
				Expand: "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations",
				ID:     "10001",
				Self:   "http://test.com/jira/rest/api/2/issue/10001",
				Key:    "ABC-123",
			},
		},
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}
