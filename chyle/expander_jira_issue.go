package chyle

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/antham/envh"
)

// JiraIssueExpander fetch data using jira issue api
type JiraIssueExpander struct {
	client   http.Client
	username string
	password string
	URL      string
	keys     map[string]string
}

// NewJiraIssueExpanderFromPasswordAuth create a new JiraIssueExpander
func NewJiraIssueExpanderFromPasswordAuth(client http.Client, username string, password string, URL string, keys map[string]string) (JiraIssueExpander, error) {
	return JiraIssueExpander{client, username, password, URL, keys}, nil
}

// Expand fetch remote jira service if a jiraIssueId is defined to fetch issue datas
func (j JiraIssueExpander) Expand(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var ID string

	if data, ok := (*commitMap)["jiraIssueId"]; true {
		if !ok {
			return commitMap, nil
		}

		if data, ok := data.(string); ok {
			ID = data
		}
	}

	req, err := http.NewRequest("GET", j.URL+"/rest/api/2/issue/"+ID, nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(j.username, j.password)
	req.Header.Set("Content-Type", "application/json")

	rep, err := j.client.Do(req)

	if err != nil {
		return commitMap, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = rep.Write(buf)

	if err != nil {
		return commitMap, err
	}

	for identifier, key := range j.keys {
		(*commitMap)[identifier] = nil

		if gjson.Get(buf.String(), key).Exists() {
			(*commitMap)[identifier] = gjson.Get(buf.String(), key).Value()
		}
	}

	return commitMap, nil
}

func buildJiraExpander(config *envh.EnvTree) (Expander, error) {
	datas := map[string]string{}
	keyValues := map[string]string{}

	for _, k := range []string{"USERNAME", "PASSWORD", "URL"} {
		v, err := config.FindString("CREDENTIALS", k)

		if err != nil {
			return nil, fmt.Errorf(`"%s" variable not found in "JIRA" config`, k)
		}

		datas[k] = v
	}

	_, err := url.ParseRequestURI(datas["URL"])

	if err != nil {
		return nil, fmt.Errorf(`"%s" is not a valid absolute URL defined in "JIRA" config`, datas["URL"])
	}

	debug(`Expander "USERNAME" defined with value "%s"`, datas["USERNAME"])
	debug(`Expander "PASSWORD" defined`)
	debug(`Expander "URL" defined with value "%s"`, datas["URL"])

	keys, err := config.FindChildrenKeys("KEYS")

	if err != nil {
		return nil, fmt.Errorf(`No "EXPANDERS_JIRA_KEYS" key found`)
	}

	for _, k := range keys {
		key, err := config.FindString("KEYS", k, "DESTKEY")

		if err != nil {
			return nil, fmt.Errorf(`An environment variable suffixed with "DESTKEY" must be defined with "%s", like EXPANDERS_JIRA_KEYS_%s_DESTKEY`, k, k)
		}

		value, err := config.FindString("KEYS", k, "FIELD")

		if err != nil {
			return nil, fmt.Errorf(`An environment variable suffixed with "FIELD" must be defined with "%s", like EXPANDERS_JIRA_KEYS_%s_FIELD`, k, k)
		}

		debug(`Expander KEY "%s" defined with value "%s"`, key, value)

		keyValues[key] = value
	}

	return NewJiraIssueExpanderFromPasswordAuth(http.Client{}, datas["USERNAME"], datas["PASSWORD"], datas["URL"], keyValues)
}
