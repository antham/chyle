package chyle

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"

	"github.com/antham/envh"
)

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

	j, err := NewJiraIssueExpanderFromPasswordAuth(*client, "test", "test", "http://test.com", map[string]string{"jiraIssueKey": "key"})

	assert.NoError(t, err, "Must return no errors")

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
			"test":         "test1",
			"jiraIssueId":  "10000",
			"jiraIssueKey": "EX-1",
		},
		map[string]interface{}{
			"test":         "test2",
			"jiraIssueId":  "ABC-123",
			"jiraIssueKey": "ABC-123",
		},
	}

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, expected, *result, "Must return same struct than the one submitted")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestCreateExpanders(t *testing.T) {
	restoreEnvs()
	setenv("EXPANDERS_JIRA_CREDENTIALS_USERNAME", "test")
	setenv("EXPANDERS_JIRA_CREDENTIALS_PASSWORD", "test")
	setenv("EXPANDERS_JIRA_CREDENTIALS_URL", "http://test.com")
	setenv("EXPANDERS_JIRA_KEYS_JIRATICKETDESCRIPTION_DESTKEY", "jiraTicketDescription")
	setenv("EXPANDERS_JIRA_KEYS_JIRATICKETDESCRIPTION_FIELD", "fields.summary")

	config, err := envh.NewEnvTree("^EXPANDERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("EXPANDERS")

	assert.NoError(t, err, "Must return no errors")

	r, err := CreateExpanders(&subConfig)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *r, 1, "Must return 1 expander")
}

func TestCreateExpandersWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("EXPANDERS_TEST", "")
			},
			`"TEST" is not a valid expander structure`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^EXPANDERS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("EXPANDERS")

		assert.NoError(t, err, "Must return no errors")

		_, err = CreateExpanders(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
