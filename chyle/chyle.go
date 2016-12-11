package chyle

import (
	"github.com/spf13/viper"
)

// BuildChangelog creates a changelog from defined configuration and from given
// range of commit references
func BuildChangelog(repoPath string, viper *viper.Viper, fromRef string, toRef string) error {
	commits, err := fetchCommits(repoPath, fromRef, toRef)

	if err != nil {
		return err
	}

	p, err := buildProcess(viper)

	if err != nil {
		return err
	}

	return proceed(p, commits)
}
