package chyle

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"

	"github.com/antham/envh"
)

func TestJiraDecorator(t *testing.T) {
	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j, err := newJiraIssueDecoratorFromPasswordAuth(*client, "test", "test", "http://test.com", map[string]string{"jiraIssueKey": "key"})

	assert.NoError(t, err, "Must return no errors")

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

	j, err := newJiraIssueDecoratorFromPasswordAuth(*client, "test", "test", "http://test.com", map[string]string{"jiraIssueKey": "key"})

	assert.NoError(t, err, "Must return no errors")

	result, err := j.decorate(&map[string]interface{}{"test": "test"})

	expected := map[string]interface{}{
		"test": "test",
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.False(t, gock.IsDone(), "Must have one pending request")
}

func TestCreateJiraDecoratorWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS", "test")
			},
			`"USERNAME" variable not found in "JIRA" config`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
			},
			`"PASSWORD" variable not found in "JIRA" config`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
				setenv("DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
			},
			`"URL" variable not found in "JIRA" config`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
				setenv("DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("DECORATORS_JIRA_CREDENTIALS_URL", "url")
			},
			`"url" is not a valid absolute URL defined in "JIRA" config`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
				setenv("DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
			},
			`No "DECORATORS_JIRA_KEYS" key found`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
				setenv("DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("DECORATORS_JIRA_KEYS_TEST", "test")
			},
			`An environment variable suffixed with "DESTKEY" must be defined with "TEST", like DECORATORS_JIRA_KEYS_TEST_DESTKEY`,
		},
		{
			func() {
				setenv("DECORATORS_JIRA_CREDENTIALS_USERNAME", "username")
				setenv("DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("DECORATORS_JIRA_KEYS_TEST_DESTKEY", "test")
			},
			`An environment variable suffixed with "FIELD" must be defined with "TEST", like DECORATORS_JIRA_KEYS_TEST_FIELD`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^DECORATORS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("DECORATORS")

		assert.NoError(t, err, "Must return no errors")

		_, err = createDecorators(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
