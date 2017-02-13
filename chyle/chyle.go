package chyle

import (
	"fmt"
	"strings"

	"github.com/antham/envh"
)

// EnableDebugging activates step logging
var EnableDebugging = false

// BuildChangelog creates a changelog from defined configuration
func BuildChangelog(envTree *envh.EnvTree) error {
	var repoPath string
	var fromRef string
	var toRef string

	for _, s := range []struct {
		keyChain []string
		ref      *string
	}{
		{
			[]string{"CHYLE", "GIT", "REPOSITORY", "PATH"},
			&repoPath,
		},
		{
			[]string{"CHYLE", "GIT", "REFERENCE", "FROM"},
			&fromRef,
		},
		{
			[]string{"CHYLE", "GIT", "REFERENCE", "TO"},
			&toRef,
		},
	} {
		ref, err := envTree.FindString(s.keyChain...)

		*(s.ref) = ref

		if err != nil {
			return fmt.Errorf("Check you defined %s", strings.Join(s.keyChain, "_"))
		}
	}

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
