package chyle

import (
	"github.com/andygrunwald/go-jira"
)

// Expander extends data from commit hashmap with data picked from third part service
type Expander interface {
	Expand(*map[string]interface{}) (*map[string]interface{}, error)
}

// JiraIssueExpander fetch data using jira issue api
type JiraIssueExpander struct {
	client *jira.Client
}

// NewJiraIssueExpanderFromPasswordAuth create a new JiraIssueExpander
func NewJiraIssueExpanderFromPasswordAuth(username string, password string, URL string) (JiraIssueExpander, error) {
	c, err := jira.NewClient(nil, URL)

	if err != nil {
		return JiraIssueExpander{}, err
	}

	res, err := c.Authentication.AcquireSessionCookie(username, password)

	if err != nil || res == false {
		return JiraIssueExpander{}, err
	}

	return JiraIssueExpander{c}, nil
}

// Expand fecth remote jira service if a jiraIssueId is defined to fetch issue datas
func (j JiraIssueExpander) Expand(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID string

	if data, ok := (*commitMap)["jiraIssueId"]; ok {
		if data, ok := data.(string); ok {
			ID = data
		}
	}

	if ID == "" {
		return commitMap, nil
	}

	issue, _, err := j.client.Issue.Get(ID)

	if err != nil {
		return commitMap, err
	}

	(*commitMap)["jiraIssue"] = issue

	return commitMap, nil
}

// Expand process all defined expander and apply them against every commit map
func Expand(expanders *[]Expander, commitMaps *[]map[string]interface{}) (*[]map[string]interface{}, error) {
	var err error

	results := []map[string]interface{}{}

	for _, commitMap := range *commitMaps {
		result := &commitMap

		for _, expander := range *expanders {
			result, err = expander.Expand(&commitMap)

			if err != nil {
				return nil, err
			}
		}

		results = append(results, *result)
	}

	return &results, nil
}
