package matchers

// Config centralizes config needed for each matcher to being
// used by any third part package to make matchers work
type Config map[string]string

// Features gives the informations if matchers are enabled
type Features struct {
	ENABLED bool
}
