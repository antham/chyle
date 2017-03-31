package chyle

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"srcd.works/go-git.v4/plumbing"
	"srcd.works/go-git.v4/plumbing/object"

	"github.com/antham/envh"
)

func getCommitFromRef(ref string) *object.Commit {
	cmd := exec.Command("git", "rev-parse", ref)
	cmd.Dir = "test"

	ID, err := cmd.Output()
	ID = ID[:len(ID)-1]

	if err != nil {
		logrus.WithField("ID", string(ID)).Fatal(err)
	}

	c, err := repo.Commit(plumbing.NewHash(string(ID)))

	if err != nil {
		logrus.WithField("ID", ID).Fatal(err)
	}

	return c
}

func TestMatchersMergeCommits(t *testing.T) {
	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err := object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := filter(&[]matcher{mergeCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "Merge branch 'test1' into test\n", "Must return merge commit message")
	assert.Equal(t, (*cs)[1].Message, "Merge branch 'test2' into test1\n", "Must return merge commit message")
}

func TestMatchersRegularCommits(t *testing.T) {
	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err := object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := filter(&[]matcher{regularCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 8, "Must return 8 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[7].Message, "feat(file5) : new file 5\n\ncreate a new file 5\n", "Must return a commit message")
}

func TestMatchersWithCommitMessage(t *testing.T) {
	re, err := regexp.Compile(`new\s*file\s*[7|8]`)

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err = object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := filter(&[]matcher{messageMatcher{re}}, &commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[1].Message, "feat(file7) : new file 7\n\ncreate a new file 7\n", "Must return a commit message")
}

func TestMatchersWithAuthor(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err = object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := filter(&[]matcher{authorMatcher{re}}, &commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}

func TestMatchersWithCommitter(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err = object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := filter(&[]matcher{committerMatcher{re}}, &commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}

func TestTransformCommitsToMap(t *testing.T) {
	commits := []object.Commit{}
	commit := getCommitFromRef("HEAD")

	err := object.WalkCommitHistory(commit, func(c *object.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	commitMaps := TransformCommitsToMap(&commits)

	expected := map[string]interface{}{
		"id":             commit.ID().String(),
		"authorName":     commit.Author.Name,
		"authorEmail":    commit.Author.Email,
		"authorDate":     commit.Author.When.String(),
		"committerName":  commit.Committer.Name,
		"committerEmail": commit.Committer.Email,
		"committerDate":  commit.Committer.When.String(),
		"message":        commit.Message,
		"type":           "regular",
	}

	assert.Len(t, *commitMaps, 10, "Must contains all history")
	assert.Equal(t, expected, (*commitMaps)[0], "Must return a map with some informations contained in commit")
}

func TestCreateMatchers(t *testing.T) {
	restoreEnvs()
	setenv("MATCHERS_TYPE", "regular")
	setenv("MATCHERS_MESSAGE", ".*")
	setenv("MATCHERS_AUTHOR", ".*")
	setenv("MATCHERS_COMMITTER", ".*")

	config, err := envh.NewEnvTree("^MATCHERS", "_")

	assert.NoError(t, err, "Must return no errors")

	subConfig, err := config.FindSubTree("MATCHERS")

	assert.NoError(t, err, "Must return no errors")

	r, err := createMatchers(&subConfig)

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *r, 4, "Must return 4 matchers")
}

func TestCreateMatchersWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		{
			func() {
				setenv("MATCHERS_TEST", "")
			},
			`a wrong matcher key containing "TEST" was defined`,
		},
		{
			func() {
				setenv("MATCHERS_TYPE", "test")
			},
			`"TYPE" must be "regular" or "merge", "test" given`,
		},
		{
			func() {
				setenv("MATCHERS_MESSAGE", "*")
			},
			`"MESSAGE" doesn't contain a valid regular expression`,
		},
		{
			func() {
				setenv("MATCHERS_COMMITTER", "*")
			},
			`"COMMITTER" doesn't contain a valid regular expression`,
		},
		{
			func() {
				setenv("MATCHERS_AUTHOR", "*")
			},
			`"AUTHOR" doesn't contain a valid regular expression`,
		},
	}

	for _, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^MATCHERS", "_")

		assert.NoError(t, err, "Must return no errors")

		subConfig, err := config.FindSubTree("MATCHERS")

		assert.NoError(t, err, "Must return no errors")

		_, err = createMatchers(&subConfig)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, test.e, "Must match error string")
	}
}
