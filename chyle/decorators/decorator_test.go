package decorators

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"

	"github.com/antham/chyle/chyle/types"
)

func TestDecorator(t *testing.T) {
	config := jiraIssueConfig{}
	config.CREDENTIALS.USERNAME = "test"
	config.CREDENTIALS.PASSWORD = "test"
	config.CREDENTIALS.URL = "http://test.com"
	config.KEYS = map[string]struct {
		DESTKEY string
		FIELD   string
	}{
		"KEY": {
			"jiraIssueKey",
			"key",
		},
	}

	defer gock.Off()

	gock.New("http://test.com/rest/api/2/issue/10000").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10000","self":"http://test.com/jira/rest/api/2/issue/10000","key":"EX-1","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	gock.New("http://test.com/rest/api/2/issue/ABC-123").
		Reply(200).
		BodyString(`{"expand":"renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations","id":"10001","self":"http://test.com/jira/rest/api/2/issue/10001","key":"ABC-123","names":{"watcher":"watcher","attachment":"attachment","sub-tasks":"sub-tasks","description":"description","project":"project","comment":"comment","issuelinks":"issuelinks","worklog":"worklog","updated":"updated","timetracking":"timetracking"	}}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	j := jiraIssueDecorator{*client, config}

	decorators := map[string][]Decorater{
		"datas":     {j},
		"metadatas": {},
	}

	changelog := types.Changelog{
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

	result, err := Decorate(&decorators, &changelog)

	expected := types.Changelog{
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

func TestCreateDataDecorators(t *testing.T) {
	tests := []func() (Features, Config){
		func() (Features, Config) {
			config := jiraIssueConfig{}
			config.CREDENTIALS.USERNAME = "test"
			config.CREDENTIALS.PASSWORD = "test"
			config.CREDENTIALS.URL = "http://test.com"
			config.KEYS = map[string]struct {
				DESTKEY string
				FIELD   string
			}{
				"DESCRIPTION": {
					"jiraTicketDescription",
					"fields.summary",
				},
			}

			return Features{JIRAISSUE: true}, Config{JIRAISSUE: config}
		},
		func() (Features, Config) {
			config := githubIssueConfig{}
			config.CREDENTIALS.OWNER = "test"
			config.CREDENTIALS.OAUTHTOKEN = "test"
			config.REPOSITORY.NAME = "test"
			config.KEYS = map[string]struct {
				DESTKEY string
				FIELD   string
			}{
				"DESCRIPTION": {
					"jiraTicketDescription",
					"fields.summary",
				},
			}

			return Features{GITHUBISSUE: true}, Config{GITHUBISSUE: config}
		},
	}

	for _, f := range tests {
		features, config := f()

		s := Create(features, config)

		assert.Len(t, (*s)["datas"], 1)
		assert.Len(t, (*s)["metadatas"], 0)
	}
}

func TestCreateMetadataDecorators(t *testing.T) {
	tests := []func() (Features, Config){
		func() (Features, Config) {
			config := envConfig{
				"TEST": {
					"TEST",
					"test",
				},
			}

			return Features{ENV: true}, Config{ENV: config}
		},
	}

	for _, f := range tests {
		features, config := f()

		s := Create(features, config)

		assert.Len(t, (*s)["datas"], 0)
		assert.Len(t, (*s)["metadatas"], 1)
	}
}
