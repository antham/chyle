package extractors

import (
	"regexp"
)

// Config centralizes config needed for each extractor to being
// used by any third part package to make extractors work
type Config map[string]struct {
	ORIGKEY string
	DESTKEY string
	REG     *regexp.Regexp
}

// Features gives tell if extractors are enabled
type Features struct {
	ENABLED bool
}
