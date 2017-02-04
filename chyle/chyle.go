package chyle

import (
	"github.com/antham/envh"
)

// EnableDebugging activates step logging
var EnableDebugging = false

// BuildChangelog creates a changelog from defined configuration and from given
// range of commit references
func BuildChangelog(repoPath string, envTree *envh.EnvTree, fromRef string, toRef string) error {
	commits, err := fetchCommits(repoPath, fromRef, toRef)

	if err != nil {
		return err
	}

	p, err := buildProcess(envTree)

	if err != nil {
		return err
	}

	return proceed(p, commits)
}
