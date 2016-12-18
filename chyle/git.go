package chyle

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// resolveRef give hash commit for a given string reference
func resolveRef(refCommit string, repository *git.Repository) (*object.Commit, error) {
	hash := plumbing.Hash{}

	if strings.ToLower(refCommit) == "head" {
		head, err := repository.Head()

		if err == nil {
			return repository.Commit(head.Hash())
		}
	}

	iter, err := repository.References()

	if err != nil {
		return &object.Commit{}, err
	}

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == refCommit {
			hash = ref.Hash()
		}

		return nil
	})

	if err == nil && !hash.IsZero() {
		return repository.Commit(hash)
	}

	hash = plumbing.NewHash(refCommit)

	if !hash.IsZero() {
		return repository.Commit(hash)
	}

	return &object.Commit{}, fmt.Errorf(`Can't find reference "%s"`, refCommit)
}

// fetchCommits retrieves commits between a reference range
func fetchCommits(repoPath string, fromRef string, toRef string) (*[]object.Commit, error) {
	commits := []object.Commit{}
	repo, err := git.NewFilesystemRepository(repoPath + "/.git/")

	if err != nil {
		return nil, err
	}

	fromCommit, err := resolveRef(fromRef, repo)

	if err != nil {
		return &[]object.Commit{}, err
	}

	toCommit, err := resolveRef(toRef, repo)

	if err != nil {
		return &[]object.Commit{}, err
	}

	_ = object.WalkCommitHistory(toCommit, func(c *object.Commit) error {
		commits = append(commits, *c)

		if c.ID() == fromCommit.ID() {
			return io.EOF
		}

		return nil
	})

	return &commits, nil
}
