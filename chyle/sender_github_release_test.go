package chyle

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v0"
)

func TestGithubReleaseSender(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/repos/test/test/releases").
		MatchHeader("Authorization", "token d41d8cd98f00b204e9800998ecf8427e").
		MatchHeader("Content-Type", "application/json").
		HeaderPresent("Accept").
		JSON(GithubRelease{TagName: "v1.0.0", Name: "TEST", Body: "Hello world !"}).
		Reply(201).
		JSON(`{
  "url": "https://api.github.com/repos/test/test/releases/1",
  "html_url": "https://github.com/test/test/releases/v1.0.0",
  "assets_url": "https://api.github.com/repos/test/test/releases/1/assets",
  "upload_url": "https://uploads.github.com/repos/test/test/releases/1/assets{?name,label}",
  "tarball_url": "https://api.github.com/repos/test/test/tarball/v1.0.0",
  "zipball_url": "https://api.github.com/repos/test/test/zipball/v1.0.0",
  "id": 1,
  "tag_name": "v1.0.0",
  "target_commitish": "master",
  "name": "v1.0.0",
  "body": "Description of the release",
  "draft": false,
  "prerelease": false,
  "created_at": "2013-02-27T19:35:32Z",
  "published_at": "2013-02-27T19:35:32Z",
  "author": {
    "login": "test",
    "id": 1,
    "avatar_url": "https://github.com/images/error/test_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/test",
    "html_url": "https://github.com/test",
    "followers_url": "https://api.github.com/users/test/followers",
    "following_url": "https://api.github.com/users/test/following{/other_user}",
    "gists_url": "https://api.github.com/users/test/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/test/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/test/subscriptions",
    "organizations_url": "https://api.github.com/users/test/orgs",
    "repos_url": "https://api.github.com/users/test/repos",
    "events_url": "https://api.github.com/users/test/events{/privacy}",
    "received_events_url": "https://api.github.com/users/test/received_events",
    "type": "User",
    "site_admin": false
  },
  "assets": [

  ]
}{"tag_name": "v1.0.0","target_commitish": "master","name": "v1.0.0","body": "Description of the release","draft": false,"prerelease": false}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	v := viper.New()
	v.Set("senders.github.template", "{{ range $key, $value := . }}{{$value.test}}{{ end }}")
	v.Set("senders.github.tag", "v1.0.0")
	v.Set("senders.github.name", "TEST")
	v.Set("senders.github.credentials.owner", "test")
	v.Set("senders.github.repository.name", "test")
	v.Set("senders.github.credentials.oauthtoken", "d41d8cd98f00b204e9800998ecf8427e")

	m, err := buildGithubReleaseSender(v)
	s := m.(GithubReleaseSender)
	s.client = *client

	assert.NoError(t, err, "Must return no errors")

	c := []map[string]interface{}{}
	c = append(c, map[string]interface{}{"test": "Hello world !"})

	err = s.Send(&c)

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, gock.IsDone(), "Must have no pending requests")
}
