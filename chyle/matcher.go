package chyle

import (
	"srcd.works/go-git.v4/plumbing/object"
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

// createMatchers build matchers from a config
func createMatchers() *[]matcher {
	results := []matcher{}

	for k, v := range chyleConfig.MATCHERS {
		results = append(results,
			map[string]func(string) matcher{
				"MESSAGE":   buildMessageMatcher,
				"COMMITTER": buildCommitterMatcher,
				"AUTHOR":    buildAuthorMatcher,
				"TYPE":      buildTypeMatcher,
			}[k](v))
	}

	return &results
}
