package chyle

import (
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var repo *git.Repository

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	for _, filename := range []string{"../features/init.sh", "../features/merge-commits.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	repo, err = git.NewFilesystemRepository(path + "/test/.git/")

	if err != nil {
		logrus.Fatal(err)
	}
}

func getCommitFromRef(ref string) *git.Commit {
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
	m := NewMatchers(&[]Matcher{MergeCommitMatcher{}})

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err := git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := m.Filter(&commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "Merge branch 'test1' into test\n", "Must return merge commit message")
	assert.Equal(t, (*cs)[1].Message, "Merge branch 'test2' into test1\n", "Must return merge commit message")
}

func TestMatchersRegularCommits(t *testing.T) {
	m := NewMatchers(&[]Matcher{RegularCommitMatcher{}})

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err := git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := m.Filter(&commits)

	assert.Len(t, *cs, 8, "Must return 8 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[7].Message, "feat(file5) : new file 5\n\ncreate a new file 5\n", "Must return a commit message")
}

func TestMatchersWithCommitMessage(t *testing.T) {
	re, err := regexp.Compile(`new\s*file\s*[7|8]`)

	if err != nil {
		logrus.Fatal(err)
	}

	m := NewMatchers(&[]Matcher{MessageMatcher{re}})

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := m.Filter(&commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[1].Message, "feat(file7) : new file 7\n\ncreate a new file 7\n", "Must return a commit message")
}

func TestMatchersWithAuthor(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	m := NewMatchers(&[]Matcher{AuthorMatcher{re}})

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := m.Filter(&commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}

func TestMatchersWithCommitter(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	m := NewMatchers(&[]Matcher{CommitterMatcher{re}})

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := m.Filter(&commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}
