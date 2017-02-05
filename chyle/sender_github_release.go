package chyle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/antham/envh"
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
	config map[string]string
}

// NewGithubReleaseSenderFromOAuth create a new GithubReleaseSender
func NewGithubReleaseSenderFromOAuth(client http.Client, config map[string]string) (GithubReleaseSender, error) {
	return GithubReleaseSender{client, config}, nil
}

// buildBody create a request body from commit map
func (j GithubReleaseSender) buildBody(commitMap *[]map[string]interface{}) ([]byte, error) {
	body, err := populateTemplate("github-release-template", j.config["TEMPLATE"], commitMap)

	if err != nil {
		return []byte{}, err
	}

	release := GithubRelease{
		TagName: j.config["TAG"],
		Name:    j.config["NAME"],
		Body:    body,
	}

	return json.Marshal(release)
}

// createRelease creates a release on github
func (j GithubReleaseSender) createRelease(body []byte) error {
	URL := "https://api.github.com/repos/%s/%s/releases"

	req, err := http.NewRequest("POST", fmt.Sprintf(URL, j.config["CREDENTIALS_OWNER"], j.config["REPOSITORY_NAME"]), bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+j.config["CREDENTIALS_OAUTHTOKEN"])
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	rep, err := j.client.Do(req)

	if err != nil {
		return ErrSenderFailed{fmt.Sprintf("can't create github release, %s", err.Error())}
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

	return ErrSenderFailed{fmt.Sprintf("can't create github release, %s", b)}
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
func buildGithubReleaseSender(config *envh.EnvTree) (Sender, error) {

	c := map[string]string{}

	for _, keyChain := range [][]string{
		[]string{"CREDENTIALS", "OAUTHTOKEN"},
		[]string{"CREDENTIALS", "OWNER"},
		[]string{"TAG"},
		[]string{"TEMPLATE"},
		[]string{"REPOSITORY", "NAME"},
	} {
		v, err := config.FindString(keyChain...)

		if err != nil {
			return nil, err
		}

		debug(`Sender GITHUB "%s" defined with value "%s"`, strings.Join(keyChain, `" "`), v)

		c[strings.Join(keyChain, "_")] = v
	}

	if config.IsExistingSubTree("NAME") {
		v, err := config.FindString("NAME")

		if err != nil {
			return nil, err
		}

		c["NAME"] = v
	}

	return NewGithubReleaseSenderFromOAuth(http.Client{}, c)
}
