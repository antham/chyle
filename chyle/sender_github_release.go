package chyle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// githubReleaseSender fetch data using jira issue api
type githubReleaseSender struct {
	client        http.Client
	config        githubReleaseConfig
	githubRelease githubRelease
}

// buildBody create a request body from commit map
func (j githubReleaseSender) buildBody(commitMap *[]map[string]interface{}) ([]byte, error) {
	body, err := populateTemplate("github-release-template", j.config.template, commitMap)

	if err != nil {
		return []byte{}, err
	}

	j.githubRelease.Body = body

	return json.Marshal(j.githubRelease)
}

func (j githubReleaseSender) requestWithPayload(URL string, method string, body []byte) error {
	req, err := http.NewRequest(method, fmt.Sprintf(URL, j.config.owner, j.config.repositoryName), bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+j.config.oauthToken)
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

	return fmt.Errorf(b)
}

// createRelease creates a release on github
func (j githubReleaseSender) createRelease(body []byte) error {
	err := j.requestWithPayload("https://api.github.com/repos/%s/%s/releases", "POST", body)

	if err != nil {
		return ErrSenderFailed{fmt.Sprintf("can't create github release : %s", err.Error())}
	}

	return nil
}

// Send push changelog to github release
func (j githubReleaseSender) Send(commitMap *[]map[string]interface{}) error {
	body, err := j.buildBody(commitMap)

	if err != nil {
		return err
	}

	return j.createRelease(body)
}
