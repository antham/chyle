package chyle

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/envh"
)

func TestBuildChangelog(t *testing.T) {
	restoreEnvs()
	p, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	setenv("CHYLE_GIT_REPOSITORY_PATH", p+"/testing-repository")
	setenv("CHYLE_GIT_REFERENCE_FROM", "test2")
	setenv("CHYLE_GIT_REFERENCE_TO", "head")
	setenv("CHYLE_MATCHERS_TYPE", "regular")
	setenv("CHYLE_EXTRACTORS_MESSAGE_ORIGKEY", "message")
	setenv("CHYLE_EXTRACTORS_MESSAGE_DESTKEY", "subject")
	setenv("CHYLE_EXTRACTORS_MESSAGE_REG", "(.{1,50})")
	setenv("CHYLE_SENDERS_STDOUT_FORMAT", "json")

	f, err := ioutil.TempFile(p+"/testing-repository", "test")

	if err != nil {
		logrus.Fatal(err)
	}

	config, err := envh.NewEnvTree("CHYLE", "_")

	if err != nil {
		logrus.Fatal(err)
	}

	oldStdout := os.Stdout
	os.Stdout = f

	err = BuildChangelog(&config)

	os.Stdout = oldStdout

	assert.NoError(t, err, "Must build changelog without errors")

	b, err := ioutil.ReadFile(f.Name())

	assert.NoError(t, err, "Can't read filename")

	type Data struct {
		ID             string `json:"id"`
		AuthorDate     string `json:"authorDate"`
		AuthorEmail    string `json:"authorEmail"`
		AuthorName     string `json:"authorName"`
		Type           string `json:"type"`
		CommitterEmail string `json:"committerEmail"`
		CommitterName  string `json:"committerName"`
		Message        string `json:"message"`
		Subject        string `json:"subject"`
	}

	results := struct {
		Datas     []Data            `json:"datas"`
		Metadatas map[string]string `json:"metadatas"`
	}{}

	j := json.NewDecoder(bytes.NewBuffer(b))
	err = j.Decode(&results)

	assert.NoError(t, err, "Must decode json without errors")
	assert.Len(t, results.Datas, 2, "Must contains 2 entries")
	assert.Len(t, results.Metadatas, 0, "Must contains no entries")

	expected := []map[string]string{
		{
			"authorEmail":    "whatever@example.com",
			"authorName":     "whatever",
			"committerEmail": "whatever@example.com",
			"committerName":  "whatever",
			"type":           "regular",
			"message":        "feat(file8) : new file 8\n\ncreate a new file 8\n",
			"subject":        "feat(file8) : new file 8",
		},
		{
			"authorEmail":    "whatever@example.com",
			"authorName":     "whatever",
			"committerEmail": "whatever@example.com",
			"committerName":  "whatever",
			"type":           "regular",
			"message":        "feat(file7) : new file 7\n\ncreate a new file 7\n",
			"subject":        "feat(file7) : new file 7",
		},
	}

	for i, r := range results.Datas {
		assert.Equal(t, expected[i]["authorEmail"], r.AuthorEmail)
		assert.Equal(t, expected[i]["authorName"], r.AuthorName)
		assert.Equal(t, expected[i]["type"], r.Type)
		assert.Equal(t, expected[i]["committerEmail"], r.CommitterEmail)
		assert.Equal(t, expected[i]["committerName"], r.CommitterName)
		assert.Equal(t, expected[i]["message"], r.Message)
		assert.Equal(t, expected[i]["subject"], r.Subject)
	}
}
