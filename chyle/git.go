package chyle

import (
	"fmt"
	"strings"

	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing"
	"srcd.works/go-git.v4/plumbing/object"
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
	repo, err := git.PlainOpen(repoPath)

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

	cs, errs := parseTree(toCommit, fromCommit)

	return &cs, concateErrors(&errs)
}

// parseTree recursively parse a given tree to extract commits till boundary is reached
func parseTree(commit *object.Commit, bound *object.Commit) ([]object.Commit, []error) {
	commits := []object.Commit{}
	errors := []error{}

	if commit.ID() == bound.ID() || commit.NumParents() == 0 {
		return commits, errors
	}

	commits = append(commits, *commit)

	parents := []object.Commit{}

	err := commit.Parents().ForEach(
		func(c *object.Commit) error {
			parents = append(parents, *c)

			return nil
		})

	if err != nil {
		errors = append(errors, err)
		return commits, errors
	}

	if len(parents) == 2 {
		cs, errs := parseTree(&parents[1], bound)
		errors = append(errors, errs...)
		commits = append(commits, cs...)
	}

	if len(parents) == 1 {
		cs, errs := parseTree(&parents[0], bound)
		errors = append(errors, errs...)
		commits = append(commits, cs...)
	}

	return commits, errors
}

// concatErrors transforms an array of error in one error
// by merging error message
func concateErrors(errs *[]error) error {
	if len(*errs) == 0 {
		return nil
	}

	errStr := ""

	for i, e := range *errs {
		errStr += e.Error()

		if i != len(*errs)-1 {
			errStr += ", "
		}
	}

	return fmt.Errorf(errStr)
}
