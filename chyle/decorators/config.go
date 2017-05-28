package decorators

// codebeat:disable[TOO_MANY_IVARS]

// Config centralizes config needed for each decorator to being
// used by any third part package to make decorators work
type Config struct {
	CUSTOMAPI   customAPIConfig
	GITHUBISSUE githubIssueConfig
	JIRAISSUE   jiraIssueConfig
	ENV         envConfig
	SHELL       shellConfig
}

// Features gives the informations if a decorator or several are defined
// and if so, which ones
type Features struct {
	ENABLED     bool
	CUSTOMAPI   bool
	JIRAISSUE   bool
	GITHUBISSUE bool
	ENV         bool
	SHELL       bool
}

// codebeat:enable[TOO_MANY_IVARS]
