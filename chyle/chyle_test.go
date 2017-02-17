package chyle

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"srcd.works/go-git.v4/plumbing"

	"github.com/antham/envh"
)

func TestBuildChangelog(t *testing.T) {
	p, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	setenv("CHYLE_GIT_REPOSITORY_PATH", p+"/test")
	setenv("CHYLE_GIT_REFERENCE_FROM", "test2")
	setenv("CHYLE_GIT_REFERENCE_TO", "head")
	setenv("CHYLE_MATCHERS_TYPE", "regular")
	setenv("CHYLE_EXTRACTORS_MESSAGE_ORIGKEY", "message")
	setenv("CHYLE_EXTRACTORS_MESSAGE_DESTKEY", "subject")
	setenv("CHYLE_EXTRACTORS_MESSAGE_REG", "(.{1,50})")
	setenv("CHYLE_SENDERS_STDOUT_FORMAT", "json")

	f, err := ioutil.TempFile(p+"/test", "test")

	if err != nil {
		logrus.Fatal(err)
	}

	oldStdout := os.Stdout
	os.Stdout = f

	config, err := envh.NewEnvTree("CHYLE", "_")

	if err != nil {
		logrus.Fatal(err)
	}

	err = BuildChangelog(&config)

	assert.NoError(t, err, "Must build changelog without errors")

	os.Stdout = oldStdout

	b, _ := ioutil.ReadFile(f.Name())

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

	results := []Data{}

	j := json.NewDecoder(bytes.NewBuffer(b))
	err = j.Decode(&results)

	assert.NoError(t, err, "Must decode json without errors")
	assert.Len(t, results, 2, "Must contains 2 entries")

	subjectExpected := []string{
		"feat(file8) : new file 8",
		"feat(file7) : new file 7",
	}

	for i, r := range results {
		h := plumbing.NewHash(r.ID)

		c, err := repo.Commit(h)
		assert.NoError(t, err, "Must return no errors")
		assert.Equal(t, c.Message, r.Message, "Must contains commit message")
		assert.Equal(t, c.Author.Name, r.AuthorName, "Must contains author name")
		assert.Equal(t, c.Author.Email, r.AuthorEmail, "Must contains author email")
		assert.Equal(t, c.Author.When.String(), r.AuthorDate, "Must contains commit date")
		assert.Equal(t, c.Committer.Name, r.CommitterName, "Must contains committer name")
		assert.Equal(t, c.Committer.Email, r.CommitterEmail, "Must contains committer email")
		assert.Equal(t, subjectExpected[i], r.Subject, "Must contains a subject field")
		assert.Equal(t, r.Type, "regular", "Must have a commit type")
	}
}
