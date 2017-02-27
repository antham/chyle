package chyle

import (
	"github.com/antham/envh"
)

// EnableDebugging activates step logging
var EnableDebugging = false

// changelogConfig stores base config needed to generate a changelog
type changelogConfig struct {
	path string
	from string
	to   string
}

// BuildChangelog creates a changelog from defined configuration
func BuildChangelog(config *envh.EnvTree) error {
	cConfig, err := extractChangelogConfig(config)

	if err != nil {
		return err
	}

	commits, err := fetchCommits(cConfig.path, cConfig.from, cConfig.to)

	if err != nil {
		return err
	}

	p, err := buildProcess(config)

	if err != nil {
		return err
	}

	return proceed(p, commits)
}

// extractChangelogConfig parses initial config
func extractChangelogConfig(config *envh.EnvTree) (changelogConfig, error) {
	cConfig := changelogConfig{}

	return cConfig, extractStringConfig(
		config,
		[]strConfigMapping{
			strConfigMapping{
				[]string{"CHYLE", "GIT", "REPOSITORY", "PATH"},
				&cConfig.path,
				true,
			},
			strConfigMapping{
				[]string{"CHYLE", "GIT", "REFERENCE", "FROM"},
				&cConfig.from,
				true,
			},
			strConfigMapping{
				[]string{"CHYLE", "GIT", "REFERENCE", "TO"},
				&cConfig.to,
				true,
			},
		},
		[]string{},
	)
}
