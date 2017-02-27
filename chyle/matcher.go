package chyle

import (
	"fmt"
	"srcd.works/go-git.v4/plumbing/object"

	"github.com/antham/envh"
)

// Matcher describe a way of applying a matcher against a commit
type Matcher interface {
	Match(*object.Commit) bool
}

// Filter commits that don't fit any matchers
func Filter(matchers *[]Matcher, commits *[]object.Commit) *[]object.Commit {
	results := []object.Commit{}

	for _, commit := range *commits {
		add := true
		for _, matcher := range *matchers {
			if !matcher.Match(&commit) {
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

// CreateMatchers build matchers from a config
func CreateMatchers(config *envh.EnvTree) (*[]Matcher, error) {
	results := []Matcher{}

	var m Matcher
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

			m, err = map[string]func(string, string) (Matcher, error){
				"MESSAGE":   buildMessageMatcher,
				"COMMITTER": buildCommitterMatcher,
				"AUTHOR":    buildAuthorMatcher,
				"TYPE":      buildTypeMatcher,
			}[k](k, s)
		default:
			err = fmt.Errorf(`a wrong matcher key containing "%s" was defined`, k)
		}

		if err != nil {
			return &[]Matcher{}, err
		}

		results = append(results, m)
	}

	return &results, nil
}
