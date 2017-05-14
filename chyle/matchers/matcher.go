package matchers

import (
	"srcd.works/go-git.v4/plumbing/object"
)

// Matcher describe a way of applying a matcher against a commit
type Matcher interface {
	Match(*object.Commit) bool
}

// Filter commits that don't fit any matchers
func Filter(matchers *[]Matcher, commits *[]object.Commit) *[]map[string]interface{} {
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

	return transformCommitsToMap(&results)
}

// transformCommitsToMap extract useful commits data in hash map table
func transformCommitsToMap(commits *[]object.Commit) *[]map[string]interface{} {
	var commitMap map[string]interface{}
	commitMaps := []map[string]interface{}{}

	for _, c := range *commits {
		commitMap = map[string]interface{}{
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

// Create builds matchers from a config
func Create(features Features, matchers Config) *[]Matcher {
	results := []Matcher{}

	if features.AUTHOR {
		results = append(results, buildAuthor(matchers.AUTHOR))
	}

	if features.COMMITTER {
		results = append(results, buildCommitter(matchers.COMMITTER))
	}

	if features.MESSAGE {
		results = append(results, buildMessage(matchers.MESSAGE))
	}

	if features.TYPE {
		results = append(results, buildType(matchers.TYPE))
	}

	return &results
}
