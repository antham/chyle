package chyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrGithubSender handles error regarding github api
// it outputs direct errors coming from request bu as well
// api return payload
type ErrGithubSender struct {
	msg string
	err error
}

func (e ErrGithubSender) Error() string {
	return fmt.Sprintf("%s : %s", e.msg, e.err)
}

// githubReleaseSender fetch data using jira issue api
type githubReleaseSender struct {
	client        *http.Client
	config        githubReleaseConfig
	githubRelease githubRelease
}

// newGithubReleaseSender creates a new githubReleaseSender object
func newGithubReleaseSender(client *http.Client, config githubReleaseConfig, githubRelease githubRelease) githubReleaseSender {
	return githubReleaseSender{
		client,
		config,
		githubRelease,
	}
}

// buildBody create a request body from changelog
func (j githubReleaseSender) buildBody(changelog *Changelog) ([]byte, error) {
	body, err := populateTemplate("github-release-template", j.config.template, changelog)

	if err != nil {
		return []byte{}, err
	}

	j.githubRelease.Body = body

	return json.Marshal(j.githubRelease)
}

// createRelease creates a release on github
func (j githubReleaseSender) createRelease(body []byte) error {
	errMsg := "can't create github release"

	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", j.config.owner, j.config.repositoryName)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + j.config.oauthToken,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(j.client, req)

	if err != nil {
		return ErrGithubSender{errMsg, err}
	}

	if status != 201 {
		return ErrGithubSender{errMsg, fmt.Errorf(string(body))}
	}

	return nil
}

// getReleaseID retrieves github release ID from a given tag name
func (j githubReleaseSender) getReleaseID() (int, error) {
	type s struct {
		ID int `json:"id"`
	}

	release := s{}

	errMsg := fmt.Sprintf("can't retrieve github release %s", j.githubRelease.TagName)
	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", j.config.owner, j.config.repositoryName, j.githubRelease.TagName)

	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return 0, err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + j.config.oauthToken,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(j.client, req)

	if err != nil {
		return 0, ErrGithubSender{errMsg, err}
	}

	if status != 200 {
		return 0, ErrGithubSender{errMsg, fmt.Errorf(string(body))}
	}

	err = json.Unmarshal(body, &release)

	if err != nil {
		return 0, ErrGithubSender{errMsg, fmt.Errorf("can't decode json body")}
	}

	return release.ID, nil
}

// updateRelease updates an existing release from a tag name
func (j githubReleaseSender) updateRelease(body []byte) error {
	ID, err := j.getReleaseID()

	if err != nil {
		return err
	}

	errMsg := fmt.Sprintf("can't update github release %s", j.githubRelease.TagName)
	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", j.config.owner, j.config.repositoryName, ID)

	req, err := http.NewRequest("PATCH", URL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + j.config.oauthToken,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(j.client, req)

	if err != nil {
		return ErrGithubSender{errMsg, err}
	}

	if status != 200 {
		return ErrGithubSender{errMsg, fmt.Errorf(string(body))}
	}

	return nil
}

// Send push changelog to github release
func (j githubReleaseSender) Send(changelog *Changelog) error {
	body, err := j.buildBody(changelog)

	if err != nil {
		return err
	}

	if j.config.update {
		return j.updateRelease(body)
	}

	return j.createRelease(body)
}
