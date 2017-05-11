package chyle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// codebeat:disable[TOO_MANY_IVARS]

// githubRelease follows https://developer.github.com/v3/repos/releases/#create-a-release
type githubRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	PreRelease      bool   `json:"prerelease,omitempty"`
}

// codebeat:enable[TOO_MANY_IVARS]
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

	_, _, err = sendRequest(g.client, req)

	return addCustomMessageToError("can't create github release", err)
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

	_, body, err := sendRequest(g.client, req)

	if err != nil {
		return 0, addCustomMessageToError(errMsg, err)
	}

	err = json.Unmarshal(body, &release)

	if err != nil {
		return 0, addCustomMessageToError(errMsg, err)
	}

	return release.ID, nil
}

// updateRelease updates an existing release from a tag name
func (g githubReleaseSender) updateRelease(body []byte) error {
	ID, err := g.getReleaseID()

	if err != nil {
		return err
	}

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

	_, _, err = sendRequest(g.client, req)

	return addCustomMessageToError(fmt.Sprintf("can't update github release %s", chyleConfig.SENDERS.GITHUB.RELEASE.TAGNAME), err)
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
