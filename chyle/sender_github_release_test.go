package chyle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestGithubReleaseSenderCreateRelease(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASGITHUBRELEASESENDER = true
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "{{ range $key, $value := .Datas }}{{$value.test}}{{ end }}"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "v1.0.0"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.NAME = "TEST"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "user"
	chyleConfig.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"

	defer gock.Off()

	tagCreationResponse, err := ioutil.ReadFile("fixtures/1-github-tag-creation-response.json")

	assert.NoError(t, err, "Must read json fixture file")

	gock.New("https://api.github.com").
		Post("/repos/user/test/releases").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(githubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		Reply(201).
		JSON(string(tagCreationResponse))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	s := buildGithubReleaseSender().(githubReleaseSender)
	s.client = client

	c := Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = append(c.Datas, map[string]interface{}{"test": "Hello world !"})

	err = s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubReleaseSenderCreateReleaseWithWrongCredentials(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASGITHUBRELEASESENDER = true
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "{{ range $key, $value := .Datas }}{{$value.test}}{{ end }}"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "v1.0.0"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.NAME = "TEST"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "d0b934ea223577f7e5cc6599e40b1822"

	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/test/test/releases").
		MatchHeader("Authorization", "token d0b934ea223577f7e5cc6599e40b1822").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(githubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		ReplyError(fmt.Errorf("an error occurred"))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	s := buildGithubReleaseSender().(githubReleaseSender)
	s.client = client

	c := Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = append(c.Datas, map[string]interface{}{"test": "Hello world !"})

	err := s.Send(&c)

	assert.EqualError(t, err, "can't create github release : Post https://api.github.com/repos/test/test/releases: an error occurred", "Must return an error when api response something wrong")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubReleaseUpdateReleaseSender(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASGITHUBRELEASESENDER = true
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "{{ range $key, $value := .Datas }}{{$value.test}}{{ end }}"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "v1.0.0"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.NAME = "TEST"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.UPDATE = true

	defer gock.Off()

	fetchReleaseResponse, err := ioutil.ReadFile("fixtures/2-github-release-fetch-response.json")

	assert.NoError(t, err, "Must read json fixture file")

	gock.New("https://api.github.com").
		Get("/repos/test/test/releases/tags/v1.0.0").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		Reply(200).
		JSON(string(fetchReleaseResponse))

	gock.New("https://api.github.com").
		Patch("/repos/test/test/releases/1").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(githubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		Reply(200)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	s := buildGithubReleaseSender().(githubReleaseSender)
	s.client = client

	c := Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = append(c.Datas, map[string]interface{}{"test": "Hello world !"})

	err = s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}

func TestGithubReleaseSenderUpdateReleaseWithWrongCredentials(t *testing.T) {
	chyleConfig = CHYLE{}
	chyleConfig.FEATURES.HASGITHUBRELEASESENDER = true
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "{{ range $key, $value := .Datas }}{{$value.test}}{{ end }}"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "v1.0.0"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.NAME = "TEST"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"
	chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "d0b934ea223577f7e5cc6599e40b1822"
	chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.UPDATE = true

	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/test/test/releases/tags/v1.0.0").
		MatchHeader("Authorization", "token d0b934ea223577f7e5cc6599e40b1822").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		ReplyError(fmt.Errorf("an error occurred"))

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	s := buildGithubReleaseSender().(githubReleaseSender)
	s.client = client

	c := Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = append(c.Datas, map[string]interface{}{"test": "Hello world !"})

	err := s.Send(&c)

	assert.EqualError(t, err, "can't retrieve github release v1.0.0 : Get https://api.github.com/repos/test/test/releases/tags/v1.0.0: an error occurred")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}
