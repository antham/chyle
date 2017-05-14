package senders

// Config centralizes config needed for each sender to being
// used by any third part package to make senders work
type Config struct {
	STDOUT        stdoutConfig
	GITHUBRELEASE githubReleaseConfig
}
