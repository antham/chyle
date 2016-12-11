package chyle

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func TestBuildChangelog(t *testing.T) {
	p, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	v := viper.New()
	v.SetConfigFile(p + "/../features/chyle.toml")
	err = v.ReadInConfig()

	if err != nil {
		logrus.Fatal(err)
	}

	f, err := ioutil.TempFile(p+"/test", "test")

	if err != nil {
		logrus.Fatal(err)
	}

	oldStdout := os.Stdout
	os.Stdout = f

	err = BuildChangelog(p+"/test", v, "test2", "head")

	os.Stdout = oldStdout

	b, _ := ioutil.ReadFile(f.Name())

	type Data struct {
		ID             string `json:"id"`
		AuthorDate     string `json:"authorDate"`
		AuthorEmail    string `json:"authorEmail"`
		AuthorName     string `json:"authorName"`
		IsMerge        bool   `json:"isMerge"`
		CommitterEmail string `json:"committerEmail"`
		CommitterName  string `json:"committerName"`
		Message        string `json:"message"`
		Subject        string `json:"subject"`
	}

	results := []Data{}

	j := json.NewDecoder(bytes.NewBuffer(b))
	err = j.Decode(&results)

	if err != nil {
		logrus.Fatal(err)
	}

	assert.Len(t, results, 7, "Must contains 7 entries")

	subjectExpected := []string{
		"feat(file8) : new file 8",
		"feat(file7) : new file 7",
		"feat(file2) : new file 2",
		"feat(file1) : new file 1",
		"feat(file4) : new file 4",
		"feat(file3) : new file 3",
		"feat(file6) : new file 6",
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
		assert.False(t, r.IsMerge, false)
	}
}
