package chyle

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestParseTree(t *testing.T) {
	type g struct {
		toRef   *object.Commit
		fromRef *object.Commit
		f       func([]object.Commit, []error)
	}

	tests := []g{
		g{
			getCommitFromRef("HEAD"),
			getCommitFromRef("test"),
			func(cs []object.Commit, errs []error) {
				assert.Len(t, errs, 0, "Must contains no errors")
				assert.Len(t, cs, 0, "Must contains 0 commits")
			},
		},
		g{
			getCommitFromRef("HEAD"),
			getCommitFromRef("test1"),
			func(cs []object.Commit, errs []error) {
				assert.Len(t, errs, 0, "Must contains no errors")
				assert.Len(t, cs, 3, "Must contains 3 commits")

				commitTests := []string{
					"feat(file8) : new file 8\n\ncreate a new file 8\n",
					"feat(file7) : new file 7\n\ncreate a new file 7\n",
					"Merge branch 'test1' into test\n",
				}

				for i, c := range cs {
					assert.Equal(t, commitTests[i], c.Message, "Must match message")
				}
			},
		},
		g{
			getCommitFromRef("HEAD"),
			getCommitFromRef("test2"),
			func(cs []object.Commit, errs []error) {
				assert.Len(t, errs, 0, "Must contains no errors")
				assert.Len(t, cs, 4, "Must contains 4 commits")

				commitTests := []string{
					"feat(file8) : new file 8\n\ncreate a new file 8\n",
					"feat(file7) : new file 7\n\ncreate a new file 7\n",
					"Merge branch 'test1' into test\n",
					"Merge branch 'test2' into test1\n",
				}

				for i, c := range cs {
					assert.Equal(t, commitTests[i], c.Message, "Must match message")
				}
			},
		},
		g{
			getCommitFromRef("HEAD"),
			getCommitFromRef("HEAD~4"),
			func(cs []object.Commit, errs []error) {
				assert.Len(t, errs, 0, "Must contains no errors")
				assert.Len(t, cs, 10, "Must contains 10 commits")

				commitTests := []string{
					"feat(file8) : new file 8\n\ncreate a new file 8\n",
					"feat(file7) : new file 7\n\ncreate a new file 7\n",
					"Merge branch 'test1' into test\n",
					"Merge branch 'test2' into test1\n",
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
					"feat(file2) : new file 2\n\ncreate a new file 2\n",
					"feat(file1) : new file 1\n\ncreate a new file 1\n",
				}

				for i, c := range cs {
					assert.Equal(t, commitTests[i], c.Message, "Must match message")
				}
			},
		},
	}

	for _, test := range tests {
		test.f(parseTree(test.toRef, test.fromRef))
	}
}

func TestConcatErrors(t *testing.T) {
	type g struct {
		errs *[]error
		f    func(error)
	}

	tests := []g{
		g{
			&[]error{},
			func(err error) {
				assert.NoError(t, err, "Must contains no error")
			},
		},
		g{
			&[]error{fmt.Errorf("test1")},
			func(err error) {
				assert.Error(t, err, "Must contains an error")
				assert.EqualError(t, err, "test1", "Must match error string")
			},
		},
		g{
			&[]error{fmt.Errorf("test1"), fmt.Errorf("test2")},
			func(err error) {
				assert.Error(t, err, "Must contains an error")
				assert.EqualError(t, err, "test1, test2", "Must match error string")
			},
		},
	}

	for _, test := range tests {
		err := concateErrors(test.errs)

		test.f(err)
	}
}
