package chyle

import (
	"github.com/antham/envh"
)

// EnableDebugging activates step logging
var EnableDebugging = false

// BuildChangelog creates a changelog from defined configuration
func BuildChangelog(envConfig *envh.EnvTree) error {
	err := resolveConfig(envConfig)

	if err != nil {
		return err
	}

	commits, err := fetchCommits(chyleConfig.GIT.REPOSITORY.PATH, chyleConfig.GIT.REFERENCE.FROM, chyleConfig.GIT.REFERENCE.TO)

	if err != nil {
		return err
	}

	p, err := buildProcess()

	if err != nil {
		return err
	}

	return proceed(p, commits)
}
