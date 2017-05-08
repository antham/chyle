package chyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// githubRelease follows https://developer.github.com/v3/repos/releases/#create-a-release
// codebeat:disable[TOO_MANY_IVARS]
type githubRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	PreRelease      bool   `json:"prerelease,omitempty"`
}

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

// buildGithubReleaseSender create a new GithubReleaseSender object from viper config
func buildGithubReleaseSender() sender {
	return newGithubReleaseSender(&http.Client{})
}

// githubReleaseSender fetch data using jira issue api
type githubReleaseSender struct {
	client *http.Client
}

// newGithubReleaseSender creates a new githubReleaseSender object
func newGithubReleaseSender(client *http.Client) githubReleaseSender {
	return githubReleaseSender{client}
}

// buildBody create a request body from changelog
func (g githubReleaseSender) buildBody(changelog *Changelog) ([]byte, error) {
	body, err := populateTemplate("github-release-template", chyleConfig.SENDERS.GITHUB.RELEASE.TEMPLATE, changelog)

	if err != nil {
		return []byte{}, err
	}

	r := githubRelease{
		chyleConfig.SENDERS.GITHUB.RELEASE.TAGNAME,
		chyleConfig.SENDERS.GITHUB.RELEASE.TARGETCOMMITISH,
		chyleConfig.SENDERS.GITHUB.RELEASE.NAME,
		body,
		chyleConfig.SENDERS.GITHUB.RELEASE.DRAFT,
		chyleConfig.SENDERS.GITHUB.RELEASE.PRERELEASE,
	}

	return json.Marshal(r)
}

// createRelease creates a release on github
func (g githubReleaseSender) createRelease(body []byte) error {
	errMsg := "can't create github release"

	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", chyleConfig.SENDERS.GITHUB.CREDENTIALS.OWNER, chyleConfig.SENDERS.GITHUB.REPOSITORY.NAME)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + chyleConfig.SENDERS.GITHUB.CREDENTIALS.OAUTHTOKEN,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(g.client, req)

	if err != nil {
		return ErrGithubSender{errMsg, err}
	}

	if status != 201 {
		return ErrGithubSender{errMsg, fmt.Errorf(string(body))}
	}

	return nil
}

// getReleaseID retrieves github release ID from a given tag name
func (g githubReleaseSender) getReleaseID() (int, error) {
	type s struct {
		ID int `json:"id"`
	}

	release := s{}

	errMsg := fmt.Sprintf("can't retrieve github release %s", chyleConfig.SENDERS.GITHUB.RELEASE.TAGNAME)
	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", chyleConfig.SENDERS.GITHUB.CREDENTIALS.OWNER, chyleConfig.SENDERS.GITHUB.REPOSITORY.NAME, chyleConfig.SENDERS.GITHUB.RELEASE.TAGNAME)

	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return 0, err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + chyleConfig.SENDERS.GITHUB.CREDENTIALS.OAUTHTOKEN,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(g.client, req)

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
func (g githubReleaseSender) updateRelease(body []byte) error {
	ID, err := g.getReleaseID()

	if err != nil {
		return err
	}

	errMsg := fmt.Sprintf("can't update github release %s", chyleConfig.SENDERS.GITHUB.RELEASE.TAGNAME)
	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", chyleConfig.SENDERS.GITHUB.CREDENTIALS.OWNER, chyleConfig.SENDERS.GITHUB.REPOSITORY.NAME, ID)

	req, err := http.NewRequest("PATCH", URL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	setHeaders(req, map[string]string{
		"Authorization": "token " + chyleConfig.SENDERS.GITHUB.CREDENTIALS.OAUTHTOKEN,
		"Content-Type":  "application/json",
		"Accept":        "application/vnd.github.v3+json",
	})

	status, body, err := sendRequest(g.client, req)

	if err != nil {
		return ErrGithubSender{errMsg, err}
	}

	if status != 200 {
		return ErrGithubSender{errMsg, fmt.Errorf(string(body))}
	}

	return nil
}

// Send push changelog to github release
func (g githubReleaseSender) Send(changelog *Changelog) error {
	body, err := g.buildBody(changelog)

	if err != nil {
		return err
	}

	if chyleConfig.SENDERS.GITHUB.RELEASE.UPDATE {
		return g.updateRelease(body)
	}

	return g.createRelease(body)
}
