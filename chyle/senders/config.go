package senders

// codebeat:disable[TOO_MANY_IVARS]

// Config centralizes config needed for each sender to being
// used by any third part package to make senders work
type Config struct {
	STDOUT        stdoutConfig
	GITHUBRELEASE githubReleaseConfig
	CUSTOMAPI     customAPIConfig
}

// Features gives the informations if a sender or several are defined
// and if so, which ones
type Features struct {
	ENABLED       bool
	GITHUBRELEASE bool
	STDOUT        bool
	CUSTOMAPI     bool
}

// codebeat:enable[TOO_MANY_IVARS]
