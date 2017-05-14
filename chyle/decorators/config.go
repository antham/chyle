package decorators

// Config centralizes config needed for each decorator to being
// used by any third part package to make decorators work
type Config struct {
	GITHUBISSUE githubIssueConfig
	JIRAISSUE   jiraIssueConfig
	ENV         envConfig
}
