package matchers

import (
	"regexp"
)

// Config centralizes config needed for each matcher to being
// used by any third part package to make matchers work
type Config struct {
	MESSAGE   *regexp.Regexp
	COMMITTER *regexp.Regexp
	AUTHOR    *regexp.Regexp
	TYPE      string
}

// Features gives the informations if matchers are enabled
type Features struct {
	ENABLED   bool
	MESSAGE   bool
	COMMITTER bool
	AUTHOR    bool
	TYPE      bool
}
