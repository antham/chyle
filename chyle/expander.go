package chyle

import (
	"fmt"
	"net/url"

	"github.com/andygrunwald/go-jira"
)

// Expander extends data from commit hashmap with data picked from third part service
type Expander interface {
	Expand(*map[string]interface{}) (*map[string]interface{}, error)
}

// JiraIssueExpander fetch data using jira issue api
type JiraIssueExpander struct {
	username string
	password string
	client   *jira.Client
}

// NewJiraIssueExpanderFromPasswordAuth create a new JiraIssueExpander
func NewJiraIssueExpanderFromPasswordAuth(username string, password string, URL string) (JiraIssueExpander, error) {
	c, err := jira.NewClient(nil, URL)

	if err != nil {
		return JiraIssueExpander{}, err
	}

	return JiraIssueExpander{username, password, c}, nil
}

// Authenticate acquire a new jira session cookie
func (j JiraIssueExpander) Authenticate() (bool, error) {
	res, err := j.client.Authentication.AcquireSessionCookie(j.username, j.password)

	if err != nil || !res {
		return false, err
	}

	return true, nil
}

// Expand fetch remote jira service if a jiraIssueId is defined to fetch issue datas
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

func buildJiraExpander(config map[string]string) (Expander, error) {
	var username, password, rawURL string
	var URL *url.URL
	var ok bool

	if username, ok = config["username"]; !ok {
		return nil, fmt.Errorf(`"username" must be defined in jira config`)
	}

	if password, ok = config["password"]; !ok {
		return nil, fmt.Errorf(`"password" must be defined in jira config`)
	}

	if rawURL, ok = config["url"]; !ok {
		return nil, fmt.Errorf(`"url" must be defined in jira config`)
	}

	URL, err := url.Parse(rawURL)

	if err != nil {
		return nil, fmt.Errorf(`"%s" not a valid URL defined in jira config`, rawURL)
	}

	return NewJiraIssueExpanderFromPasswordAuth(username, password, URL.String())
}

// CreateExpanders build expanders from a config
func CreateExpanders(expanders map[string]interface{}) (*[]Expander, error) {
	results := []Expander{}

	for dk, dv := range expanders {
		var ex Expander
		var err error

		e, ok := dv.(map[string]string)

		if !ok {
			return &[]Expander{}, fmt.Errorf(`expander "%s" must contains key=value string values`, dk)
		}

		switch dk {
		case "jira":
			ex, err = buildJiraExpander(e)
		default:
			err = fmt.Errorf(`"%s" is not a valid expander structure`, dk)
		}

		if err != nil {
			return &[]Expander{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
