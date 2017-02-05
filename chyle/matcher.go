package chyle

import (
	"fmt"
	"regexp"

	"srcd.works/go-git.v4/plumbing/object"

	"github.com/antham/envh"
)

// Matcher describe a way of applying a matcher against a commit
type Matcher interface {
	Match(*object.Commit) bool
}

// MergeCommitMatcher match merge commit message
type MergeCommitMatcher struct {
}

// Match is valid if commit is a merge commit
func (m MergeCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// RegularCommitMatcher match regular commit message
type RegularCommitMatcher struct {
}

// Match is valid if commit is not a merge commit
func (r RegularCommitMatcher) Match(commit *object.Commit) bool {
	return commit.NumParents() == 1 || commit.NumParents() == 0
}

// MessageMatcher is commit message matcher
type MessageMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit message
func (m MessageMatcher) Match(commit *object.Commit) bool {
	return m.regexp.MatchString(commit.Message)
}

// CommitterMatcher is commit committer matcher
type CommitterMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit committer field
func (c CommitterMatcher) Match(commit *object.Commit) bool {
	return c.regexp.MatchString(commit.Committer.String())
}

// AuthorMatcher is commit author matcher
type AuthorMatcher struct {
	regexp *regexp.Regexp
}

// Match apply a regexp against commit author field
func (a AuthorMatcher) Match(commit *object.Commit) bool {
	return a.regexp.MatchString(commit.Author.String())
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
			"message":        c.Message,
			"isMerge":        c.NumParents() == 2,
		}

		commitMaps = append(commitMaps, commitMap)
	}

	return &commitMaps
}

func buildNumParentsMatcher(value int) (Matcher, error) {
	switch value {
	case 1, 0:
		return RegularCommitMatcher{}, nil
	case 2:
		return MergeCommitMatcher{}, nil
	}

	return nil, fmt.Errorf(`"NUMPARENTS" must be 0, 1 or 2, "%d" given`, value)
}

func buildMessageMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return MessageMatcher{r}, nil
}

func buildCommitterMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return CommitterMatcher{r}, nil
}

func buildAuthorMatcher(key string, value string) (Matcher, error) {
	r, err := regexp.Compile(value)

	if err != nil {
		return nil, fmt.Errorf(`"%s" doesn't contain a valid regular expression`, key)
	}

	return AuthorMatcher{r}, nil
}

// CreateMatchers build matchers from a config
func CreateMatchers(config *envh.EnvTree) (*[]Matcher, error) {
	results := []Matcher{}

	var m Matcher
	var i int
	var s string
	var err error

	for _, k := range config.GetChildrenKeys() {
		switch k {
		case "NUMPARENTS":
			i, err = config.FindInt(k)

			if err != nil {
				break
			}

			debug(`Matcher "%s" defined with value "%d"`, k, i)

			m, err = buildNumParentsMatcher(i)
		case "MESSAGE", "COMMITTER", "AUTHOR":
			s, err = config.FindString(k)

			if err != nil {
				break
			}

			debug(`Matcher "%s" defined with value "%s"`, k, s)

			m, err = map[string]func(string, string) (Matcher, error){
				"MESSAGE":   buildMessageMatcher,
				"COMMITTER": buildCommitterMatcher,
				"AUTHOR":    buildAuthorMatcher,
			}[k](k, s)
		default:
			err = fmt.Errorf(`"%s" is not a valid matcher structure`, k)
		}

		if err != nil {
			return &[]Matcher{}, err
		}

		results = append(results, m)
	}

	return &results, nil
}
