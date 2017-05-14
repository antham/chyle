package chyle

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

	err := Send(&[]sender{s}, &c)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, `{"datas":[{"id":1,"test":"test"},{"id":2,"test":"test"}],"metadatas":{}}`, strings.TrimRight(buf.String(), "\n"), "Must output all commit informations  as json")
}
func TestCreateSenders(t *testing.T) {
	tests := []func(){
		func() {
			chyleConfig.FEATURES.HASSTDOUTSENDER = true
			chyleConfig.SENDERS.STDOUT.FORMAT = "json"
		},
		func() {
			chyleConfig.FEATURES.HASGITHUBRELEASESENDER = true
			chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "test"
			chyleConfig.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "test"
			chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "test"
			chyleConfig.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "test"
			chyleConfig.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"
		},
	}

	for _, f := range tests {
		chyleConfig = CHYLE{}

		f()

		s := createSenders()

		assert.Len(t, *s, 1, "Must return 1 sender")
	}
}
