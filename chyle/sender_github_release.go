package chyle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// GithubRelease follows https://developer.github.com/v3/repos/releases/#create-a-release
type GithubRelease struct {
	TagName          string `json:"tag_name"`
	TargetCommittish string `json:"target_commitish,omitempty"`
	Name             string `json:"name"`
	Body             string `json:"body"`
	Draft            bool   `json:"draft,omitempty"`
	PreRelease       bool   `json:"prerelease,omitempty"`
}

// GithubReleaseSender fetch data using jira issue api
type GithubReleaseSender struct {
	client http.Client
	config *viper.Viper
}

// NewGithubReleaseSenderFromOAuth create a new GithubReleaseSender
func NewGithubReleaseSenderFromOAuth(client http.Client, config *viper.Viper) (GithubReleaseSender, error) {
	return GithubReleaseSender{client, config}, nil
}

// buildBody create a request body from commit map
func (j GithubReleaseSender) buildBody(commitMap *[]map[string]interface{}) ([]byte, error) {
	body, err := populateTemplate("github-release-template", j.config.GetString("senders.github.template"), commitMap)

	if err != nil {
		return []byte{}, err
	}

	release := GithubRelease{
		TagName: j.config.GetString("senders.github.tag"),
		Name:    j.config.GetString("senders.github.name"),
		Body:    body,
	}

	return json.Marshal(release)
}

// createRelease creates a release on github
func (j GithubReleaseSender) createRelease(body []byte) error {
	URL := "https://api.github.com/repos/%s/%s/releases"
	owner, repositoryName, token := j.config.GetString("senders.github.credentials.owner"), j.config.GetString("senders.github.repository.name"), j.config.GetString("senders.github.credentials.oauthtoken")

	req, err := http.NewRequest("POST", fmt.Sprintf(URL, owner, repositoryName), bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	rep, err := j.client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		err = rep.Body.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	if rep.StatusCode == 201 {
		return nil
	}

	b, err := bufio.NewReader(rep.Body).ReadString('\n')

	if err != nil && err != io.EOF {
		b = "can't fetch github response"
	}

	return fmt.Errorf("Can't create github release : %s", b)
}

// Send push changelog to github release
func (j GithubReleaseSender) Send(commitMap *[]map[string]interface{}) error {
	body, err := j.buildBody(commitMap)

	if err != nil {
		return err
	}

	return j.createRelease(body)
}

// buildGithubReleaseSender create a new GithubReleaseSender object from viper config
func buildGithubReleaseSender(config *viper.Viper) (Sender, error) {
	err := checkArguments([]string{
		"senders.github.credentials.oauthtoken",
		"senders.github.credentials.owner",
		"senders.github.tag",
		"senders.github.template",
	}, config)

	if err != nil {
		return nil, err
	}

	return NewGithubReleaseSenderFromOAuth(http.Client{}, config)
}
