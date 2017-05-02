package chyle

import (
	"fmt"
	"testing"

	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"
)

func TestResolveConfig(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		{
			func() {
			},
			`environment variable missing : "CHYLE_GIT_REPOSITORY_PATH"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
			},
			`environments variables missing : "CHYLE_GIT_REFERENCE_FROM", "CHYLE_GIT_REFERENCE_TO"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
			},
			`environment variable missing : "CHYLE_GIT_REFERENCE_TO"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_MATCHERS_MESSAGE", ".**")
			},
			`provide a valid regexp for "CHYLE_MATCHERS_MESSAGE", ".**" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_MATCHERS_COMMITTER", ".**")
			},
			`provide a valid regexp for "CHYLE_MATCHERS_COMMITTER", ".**" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_MATCHERS_AUTHOR", ".**")
			},
			`provide a valid regexp for "CHYLE_MATCHERS_AUTHOR", ".**" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_MATCHERS_TYPE", "test")
			},
			`provide a value for "CHYLE_MATCHERS_TYPE" from one of those values : ["regular", "merge"], "test" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_ENV_TEST", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_ENV_TEST_DESTKEY", "CHYLE_DECORATORS_ENV_TEST_VARNAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_ENV_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_ENV_TEST_VARNAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_ENV_TEST_VARNAME", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_ENV_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
			},
			`environments variables missing : "CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "testcom")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
			},
			`provide a valid URL for "CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "testcom" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
			},
			`define at least one environment variable couple "CHYLE_DECORATORS_JIRA_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_JIRA_KEYS_*_FIELD", replace "*" with your own naming`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_TEST", "test")
			},
			`environments variables missing : "CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environments variables missing : "CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME", "CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME", "v2.0.0")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE", "{{.}}")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_TEST", "test")
			},
			`environment variable missing : "CHYLE_SENDERS_STDOUT_FORMAT"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_FORMAT", "test")
			},
			`"CHYLE_SENDERS_STDOUT_FORMAT" "test" doesn't exist`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_FORMAT", "template")
			},
			`environment variable missing : "CHYLE_SENDERS_STDOUT_TEMPLATE"`,
		},
	}
	for i, test := range tests {
		restoreEnvs()
		test.f()

		config, err := envh.NewEnvTree("^CHYLE", "_")

		assert.NoError(t, err)

		err = resolveConfig(&config)

		errDetail := fmt.Sprintf("Test %d failed", i+1)

		assert.Error(t, err, errDetail)
		assert.EqualError(t, err, test.e, errDetail)
	}
}
