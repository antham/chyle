package senders

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/chyle/chyle/types"
)

func TestSend(t *testing.T) {
	buf := &bytes.Buffer{}

	s := jSONStdoutSender{buf}

	c := types.Changelog{
		Datas:     []map[string]interface{}{},
		Metadatas: map[string]interface{}{},
	}

	c.Datas = []map[string]interface{}{
		{
			"id":   1,
			"test": "test",
		},
		{
			"id":   2,
			"test": "test",
		},
	}

	err := Send(&[]Sender{s}, &c)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `{"datas":[{"id":1,"test":"test"},{"id":2,"test":"test"}],"metadatas":{}}`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations  as json")
}
func TestCreateSenders(t *testing.T) {
	tests := []func() (map[string]bool, Config){
		func() (map[string]bool, Config) {
			config := stdoutConfig{}
			config.FORMAT = "json"

			return map[string]bool{"stdoutSender": true}, Config{STDOUT: config}
		},
		func() (map[string]bool, Config) {
			config := githubReleaseConfig{}
			config.CREDENTIALS.OAUTHTOKEN = "test"
			config.CREDENTIALS.OWNER = "test"
			config.RELEASE.TAGNAME = "test"
			config.RELEASE.TEMPLATE = "test"
			config.REPOSITORY.NAME = "test"

			return map[string]bool{"githubReleaseSender": true}, Config{GITHUBRELEASE: config}
		},
	}

	for _, f := range tests {
		features, config := f()

		s := CreateSenders(features, config)

		assert.Len(t, *s, 1, "Must return 1 sender")
	}
}
