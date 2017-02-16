package chyle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"

	"github.com/antham/envh"
)

func createTestGithubReleaseSender(t *testing.T) GithubReleaseSender {

	config, err := envh.NewEnvTree("SENDERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("SENDERS", "GITHUB")

	assert.NoError(t, err, "Must return no errors")

	releaser, err := buildGithubReleaseSender(&subConfig)

	assert.NoError(t, err, "Must return no errors")

	return releaser.(GithubReleaseSender)
}

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
	setenv("SENDERS_GITHUB_TAGNAME", "v1.0.0")
	setenv("SENDERS_GITHUB_NAME", "TEST")
	setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
	setenv("SENDERS_GITHUB_REPOSITORY_NAME", "test")
	setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")

	s := createTestGithubReleaseSender(t)
	s.client = *client

	c := []map[string]interface{}{}
	c = append(c, map[string]interface{}{"test": "Hello world !"})

	err = s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubReleaseSenderWithWrongCredentials(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/test/test/releases").
		MatchHeader("Authorization", "token d0b934ea223577f7e5cc6599e40b1822").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(GithubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		ReplyError(fmt.Errorf("an error occured"))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	restoreEnvs()
	setenv("SENDERS_GITHUB_TEMPLATE", "{{ range $key, $value := . }}{{$value.test}}{{ end }}")
	setenv("SENDERS_GITHUB_TAGNAME", "v1.0.0")
	setenv("SENDERS_GITHUB_NAME", "TEST")
	setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
	setenv("SENDERS_GITHUB_REPOSITORY_NAME", "test")
	setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d0b934ea223577f7e5cc6599e40b1822")

	s := createTestGithubReleaseSender(t)
	s.client = *client

	c := []map[string]interface{}{}
	c = append(c, map[string]interface{}{"test": "Hello world !"})

	err := s.Send(&c)

	assert.EqualError(t, err, "sender issue : can't create github release, Post https://api.github.com/repos/test/test/releases: an error occured", "Must return an error when api response something wrong")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubReleaseSenderBuildBody(t *testing.T) {
	client := http.Client{Transport: &http.Transport{}}

	s, err := NewGithubReleaseSenderFromOAuth(client, map[string]string{"TEMPLATE": "{{TEST}}}"})

	assert.NoError(t, err, "Must return no errors")

	c := []map[string]interface{}{}
	c = append(c, map[string]interface{}{"test": "Hello world !"})

	datas, err := s.buildBody(&c)

	assert.Empty(t, datas, "Must return no datas")
	assert.EqualError(t, err, `check your template is well-formed : template: github-release-template:1: function "TEST" not defined`, "Must return a template error")
}

func TestGithubReleaseSendersWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		g{
			func() {
				setenv("SENDERS_GITHUB", "test")
			},
			`missing "SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN"`,
		},
		g{
			func() {
				setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "test")
			},
			`missing "SENDERS_GITHUB_CREDENTIALS_OWNER"`,
		},
		g{
			func() {
				setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "test")
				setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
			},
			`missing "SENDERS_GITHUB_TAGNAME"`,
		},
		g{
			func() {
				setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "test")
				setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("SENDERS_GITHUB_TAGNAME", "test")
			},
			`missing "SENDERS_GITHUB_TEMPLATE"`,
		},
		g{
			func() {
				setenv("SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "test")
				setenv("SENDERS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("SENDERS_GITHUB_TAGNAME", "test")
				setenv("SENDERS_GITHUB_TEMPLATE", "test")

			},
			`missing "SENDERS_GITHUB_REPOSITORY_NAME"`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^SENDERS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("SENDERS", "GITHUB")

		assert.NoError(t, err, "Must return no errors")

		_, err = buildGithubReleaseSender(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
