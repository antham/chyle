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
