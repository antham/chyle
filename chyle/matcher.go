package chyle

import (
	"fmt"
	"srcd.works/go-git.v4/plumbing/object"

	"github.com/antham/envh"
)

// matcher describe a way of applying a matcher against a commit
type matcher interface {
	match(*object.Commit) bool
}

// filter commits that don't fit any matchers
func filter(matchers *[]matcher, commits *[]object.Commit) *[]object.Commit {
	results := []object.Commit{}

	for _, commit := range *commits {
		add := true
		for _, matcher := range *matchers {
			if !matcher.match(&commit) {
				add = false
			}
		}

		if add {
			results = append(results, commit)
		}
	}

	return &results
}

// TransformCommitsToMap extract useful commits data in hash map table
func TransformCommitsToMap(commits *[]object.Commit) *[]map[string]interface{} {
	commitMaps := []map[string]interface{}{}

	for _, c := range *commits {
		commitMap := map[string]interface{}{
			"id":             c.ID().String(),
			"authorName":     c.Author.Name,
			"authorEmail":    c.Author.Email,
			"authorDate":     c.Author.When.String(),
			"committerName":  c.Committer.Name,
			"committerEmail": c.Committer.Email,
			"committerDate":  c.Committer.When.String(),
			"message":        removePGPKey(c.Message),
			"type":           solveType(&c),
		}

		commitMaps = append(commitMaps, commitMap)
	}

	return &commitMaps
}

// createMatchers build matchers from a config
func createMatchers(config *envh.EnvTree) (*[]matcher, error) {
	results := []matcher{}

	var m matcher
	var s string
	var err error

	for _, k := range config.GetChildrenKeys() {
		switch k {
		case "MESSAGE", "COMMITTER", "AUTHOR", "TYPE":
			s, err = config.FindString(k)

			if err != nil {
				break
			}

			debug(`Matcher "%s" defined with value "%s"`, k, s)

			m, err = map[string]func(string, string) (matcher, error){
				"MESSAGE":   buildMessageMatcher,
				"COMMITTER": buildCommitterMatcher,
				"AUTHOR":    buildauthorMatcher,
				"TYPE":      buildTypeMatcher,
			}[k](k, s)
		default:
			err = fmt.Errorf(`a wrong matcher key containing "%s" was defined`, k)
		}

		if err != nil {
			return &[]matcher{}, err
		}

		results = append(results, m)
	}

	return &results, nil
}
