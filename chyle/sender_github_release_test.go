package chyle

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"

	"github.com/antham/envh"
)

func TestGithubReleaseSender(t *testing.T) {
	defer gock.Off()

	tagCreationResponse, err := ioutil.ReadFile("fixtures/1-github-tag-creation-response.json")

	assert.NoError(t, err, "Must read json fixture file")

	gock.New("https://api.github.com").
		Post("/repos/test/test/releases").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(GithubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		Reply(201).
		JSON(string(tagCreationResponse))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	restoreEnvs()
	setenv("SENDERS_GITHUB_TEMPLATE", "{{ range $key, $value := . }}{{$value.test}}{{ end }}")
	setenv("SENDERS_GITHUB_TAG", "v1.0.0")
	setenv("SENDERS_GITHUB_NAME", "TEST")
	setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
	setenv("SENDERS_GITHUB_REPOSITORY_NAME", "test")
	setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")

	config, err := envh.NewEnvTree("SENDERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("SENDERS", "GITHUB")

	assert.NoError(t, err, "Must return no errors")

	m, err := buildGithubReleaseSender(&subConfig)

	assert.NoError(t, err, "Must return no errors")

	s := m.(GithubReleaseSender)
	s.client = *client

	assert.NoError(t, err, "Must return no errors")

	c := []map[string]interface{}{}
	c = append(c, map[string]interface{}{"test": "Hello world !"})

	err = s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}
