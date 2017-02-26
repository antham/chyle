package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing/object"
)

func TestResolveRef(t *testing.T) {
	type g struct {
		ref string
		f   func(*object.Commit, error)
	}

	tests := []g{
		g{
			"HEAD",
			func(o *object.Commit, err error) {
				assert.NoError(t, err, "Must contains no error")
				assert.True(t, o.ID().String() == getCommitFromRef("HEAD").ID().String(), "Must resolve HEAD reference")
			},
		},
		g{
			"test1",
			func(o *object.Commit, err error) {
				assert.NoError(t, err, "Must contains no error")
				assert.True(t, o.ID().String() == getCommitFromRef("test1").ID().String(), "Must resolve branch reference")
			},
		},
		g{
			getCommitFromRef("test1").ID().String(),
			func(o *object.Commit, err error) {
				assert.NoError(t, err, "Must contains no error")
				assert.True(t, o.ID().String() == getCommitFromRef("test1").ID().String(), "Must resolve commit id")
			},
		},
	}

	for _, test := range tests {
		test.f(resolveRef(test.ref, repo))
	}
}

func TestResolveRefWithErrors(t *testing.T) {
	type g struct {
		ref  string
		repo *git.Repository
		f    func(*object.Commit, error)
	}

	tests := []g{
		g{
			"whatever",
			repo,
			func(o *object.Commit, err error) {
				assert.Error(t, err, "Must return an error")
				assert.EqualError(t, err, `Can't find reference "whatever"`, "Must return an error")
			},
		},
	}

	for _, test := range tests {
		test.f(resolveRef(test.ref, test.repo))
	}
}

func TestFetchCommits(t *testing.T) {
	type g struct {
		toRef   string
		fromRef string
		f       func(*[]object.Commit, error)
	}

	tests := []g{
		g{
			getCommitFromRef("HEAD").ID().String(),
			getCommitFromRef("test").ID().String(),
			func(cs *[]object.Commit, err error) {
				assert.Error(t, err, "Must return an error")
			},
		},
		g{
			getCommitFromRef("HEAD~1").ID().String(),
			getCommitFromRef("HEAD~3").ID().String(),
			func(cs *[]object.Commit, err error) {
				assert.Error(t, err, "Must return an error")
			},
		},
		g{
			getCommitFromRef("HEAD~3").ID().String(),
			getCommitFromRef("test~2^2").ID().String(),
			func(cs *[]object.Commit, err error) {
				assert.NoError(t, err, "Must return no errors")
				assert.Len(t, *cs, 5, "Must contains 3 commits")

				commitTests := []string{
					"Merge branch 'test2' into test1\n",
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
				}

				for i, c := range *cs {
					assert.Equal(t, commitTests[i], c.Message, "Must match message")
				}
			},
		},
		g{
			getCommitFromRef("HEAD~4").ID().String(),
			getCommitFromRef("test~2^2^2").ID().String(),
			func(cs *[]object.Commit, err error) {
				assert.NoError(t, err, "Must return no errors")
				assert.Len(t, *cs, 5, "Must contains 3 commits")

				commitTests := []string{
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
					"feat(file2) : new file 2\n\ncreate a new file 2\n",
				}

				for i, c := range *cs {
					assert.Equal(t, commitTests[i], c.Message, "Must match message")
				}
			},
		},
	}

	for _, test := range tests {
		test.f(fetchCommits("test", test.toRef, test.fromRef))
	}
}
