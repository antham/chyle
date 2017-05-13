package chyle

import (
	"github.com/antham/chyle/chyle/git"

	"github.com/antham/envh"
)

// EnableDebugging activates step logging
var EnableDebugging = false

// BuildChangelog creates a changelog from defined configuration
func BuildChangelog(envConfig *envh.EnvTree) error {
	if err := resolveConfig(envConfig); err != nil {
		return err
	}

	debugConfig()

	commits, err := git.FetchCommits(chyleConfig.GIT.REPOSITORY.PATH, chyleConfig.GIT.REFERENCE.FROM, chyleConfig.GIT.REFERENCE.TO)

	if err != nil {
		return err
	}

	return proceed(buildProcess(), commits)
}
