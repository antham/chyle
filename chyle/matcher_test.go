package chyle

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

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
	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err := git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := Filter(&[]Matcher{MergeCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "Merge branch 'test1' into test\n", "Must return merge commit message")
	assert.Equal(t, (*cs)[1].Message, "Merge branch 'test2' into test1\n", "Must return merge commit message")
}

func TestMatchersRegularCommits(t *testing.T) {
	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err := git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := Filter(&[]Matcher{RegularCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 8, "Must return 8 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[7].Message, "feat(file5) : new file 5\n\ncreate a new file 5\n", "Must return a commit message")
}

func TestMatchersWithCommitMessage(t *testing.T) {
	re, err := regexp.Compile(`new\s*file\s*[7|8]`)

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := Filter(&[]Matcher{MessageMatcher{re}}, &commits)

	assert.Len(t, *cs, 2, "Must return 2 objects")
	assert.Equal(t, (*cs)[0].Message, "feat(file8) : new file 8\n\ncreate a new file 8\n", "Must return a commit message")
	assert.Equal(t, (*cs)[1].Message, "feat(file7) : new file 7\n\ncreate a new file 7\n", "Must return a commit message")
}

func TestMatchersWithAuthor(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := Filter(&[]Matcher{AuthorMatcher{re}}, &commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}

func TestMatchersWithCommitter(t *testing.T) {
	re, err := regexp.Compile("@")

	if err != nil {
		logrus.Fatal(err)
	}

	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err = git.WalkCommitHistory(commit, func(c *git.Commit) error {
		commits = append(commits, *c)

		return nil
	})

	if err != nil {
		logrus.Fatal(err)
	}

	cs := Filter(&[]Matcher{CommitterMatcher{re}}, &commits)

	assert.Len(t, *cs, 10, "Must return 10 objects")
}

func TestTransformCommitsToMap(t *testing.T) {
	commits := []git.Commit{}
	commit := getCommitFromRef("HEAD")

	err := git.WalkCommitHistory(commit, func(c *git.Commit) error {
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
		"isMerge":        false,
	}

	assert.Len(t, *commitMaps, 10, "Must contains all history")
	assert.Equal(t, expected, (*commitMaps)[0], "Must return a map with some informations contained in commit")
}

func TestCreateMatchers(t *testing.T) {
	r, err := CreateMatchers(map[string]string{
		"numParents": "1",
		"message":    ".*",
		"author":     ".*",
		"committer":  ".*",
	})

	assert.NoError(t, err, "Must contains no errors")
	assert.Len(t, *r, 4, "Must return 4 matchers")
}

func TestCreateMatchersWithErrors(t *testing.T) {
	type g struct {
		s map[string]string
		e string
	}

	datas := []g{
		g{
			map[string]string{"numParents": "whatever"},
			`"numParent" is not an integer`,
		},
		g{
			map[string]string{"numParents": "3"},
			`"numParent" must be 0, 1 or 2, "3" given`,
		},
		g{
			map[string]string{"message": "**"},
			`"message" doesn't contain a valid regular expression`,
		},
		g{
			map[string]string{"committer": "**"},
			`"committer" doesn't contain a valid regular expression`,
		},
		g{
			map[string]string{"author": "**"},
			`"author" doesn't contain a valid regular expression`,
		},
		g{
			map[string]string{"whatever": "**"},
			`"whatever" is not a valid matcher structure`,
		},
	}

	for _, d := range datas {
		_, err := CreateMatchers(d.s)

		assert.Error(t, err, "Must contains an error")
		assert.EqualError(t, err, d.e, "Must match error string")
	}
}
