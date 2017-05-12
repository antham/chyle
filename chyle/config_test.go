package chyle

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"
)

func TestResolveConfigWithErrors(t *testing.T) {
	type g struct {
		f func()
		e string
	}

	tests := []g{
		// Mandatory parameters
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
		// Matchers
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
		// Extractors
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_TEST_ORIGKEY", "CHYLE_EXTRACTORS_TEST_DESTKEY", "CHYLE_EXTRACTORS_TEST_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_TEST_DESTKEY", "CHYLE_EXTRACTORS_TEST_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_DESTKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_TEST_ORIGKEY", "CHYLE_EXTRACTORS_TEST_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_REG", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_TEST_ORIGKEY", "CHYLE_EXTRACTORS_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_TEST_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_TEST_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_TEST_REG", ".**")
			},
			`provide a valid regexp for "CHYLE_EXTRACTORS_TEST_REG", ".**" given`,
		},
		// Decorators env
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
		// Decorator jira
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
		// Decorator github
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environments variables missing : "CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
			},
			`define at least one environment variable couple "CHYLE_DECORATORS_GITHUB_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_GITHUB_KEYS_*_FIELD", replace "*" with your own naming`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		// Sender github
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
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE", "{{.....}}")
			},
			`provide a valid template string for "CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE" : "template: test:1: unexpected <.> in operand", "{{.....}}" given`,
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
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE", "{{.}}")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME", "v2.0.0")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUB_REPOSITORY_NAME"`,
		},
		// Sender stdout
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
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_FORMAT", "template")
				setenv("CHYLE_SENDERS_STDOUT_TEMPLATE", "{{.....}}")
			},
			`provide a valid template string for "CHYLE_SENDERS_STDOUT_TEMPLATE" : "template: test:1: unexpected <.> in operand", "{{.....}}" given`,
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

func TestResolveConfig(t *testing.T) {
	type g struct {
		f func()
		c func() CHYLE
	}

	tests := []g{
		// Mandatory parameters
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"

				return c
			},
		},
		// Matchers
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_MATCHERS_TYPE", "regular")
				setenv("CHYLE_MATCHERS_AUTHOR", ".*")
				setenv("CHYLE_MATCHERS_COMMITTER", ".*")
				setenv("CHYLE_MATCHERS_MESSAGE", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASMATCHERS = true
				c.MATCHERS = map[string]string{
					"TYPE":      "regular",
					"AUTHOR":    ".*",
					"COMMITTER": ".*",
					"MESSAGE":   ".*",
				}

				return c
			},
		},
		// Extractors
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_EXTRACTORS_TEST_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_TEST_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_TEST_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASEXTRACTORS = true
				c.EXTRACTORS = map[string]map[string]string{
					"TEST": {
						"ORIGKEY": "test",
						"DESTKEY": "test",
						"REG":     ".*",
					},
				}

				return c
			},
		},
		// Decorators env
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_ENV_TEST_VARNAME", "var")
				setenv("CHYLE_DECORATORS_ENV_TEST_DESTKEY", "destkey")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASDECORATORS = true
				c.FEATURES.HASENVDECORATOR = true
				c.DECORATORS.ENV = map[string]map[string]string{
					"TEST": {
						"DESTKEY": "destkey",
						"VARNAME": "var",
					},
				}

				return c
			},
		},
		// Decorator jira
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRA_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_DESTKEY", "destkey")
				setenv("CHYLE_DECORATORS_JIRA_KEYS_TEST_FIELD", "field")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASEXTRACTORS = true
				c.FEATURES.HASDECORATORS = true
				c.FEATURES.HASJIRADECORATOR = true
				c.EXTRACTORS = map[string]map[string]string{
					"JIRAISSUEID": {
						"ORIGKEY": "test",
						"DESTKEY": "test",
						"REG":     ".*",
					},
				}
				c.DECORATORS.JIRA.CREDENTIALS.URL = "http://test.com"
				c.DECORATORS.JIRA.CREDENTIALS.USERNAME = "test"
				c.DECORATORS.JIRA.CREDENTIALS.PASSWORD = "password"
				c.DECORATORS.JIRA.KEYS = map[string]string{
					"destkey": "field",
				}

				return c
			},
		},
		// Decorator github
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUB_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUB_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_DESTKEY", "destkey")
				setenv("CHYLE_DECORATORS_GITHUB_KEYS_TEST_FIELD", "field")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASEXTRACTORS = true
				c.FEATURES.HASDECORATORS = true
				c.FEATURES.HASGITHUBDECORATOR = true
				c.EXTRACTORS = map[string]map[string]string{
					"GITHUBISSUEID": {
						"ORIGKEY": "test",
						"DESTKEY": "test",
						"REG":     ".*",
					},
				}
				c.DECORATORS.GITHUB.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
				c.DECORATORS.GITHUB.CREDENTIALS.OWNER = "test"
				c.DECORATORS.GITHUB.REPOSITORY.NAME = "test"
				c.DECORATORS.GITHUB.KEYS = map[string]string{
					"destkey": "field",
				}

				return c
			},
		},
		// Sender github
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUB_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TEMPLATE", "{{.}}")
				setenv("CHYLE_SENDERS_GITHUB_RELEASE_TAGNAME", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUB_REPOSITORY_NAME", "test")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASSENDERS = true
				c.FEATURES.HASGITHUBRELEASESENDER = true
				c.SENDERS.GITHUB.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
				c.SENDERS.GITHUB.CREDENTIALS.OWNER = "user"
				c.SENDERS.GITHUB.RELEASE.TAGNAME = "v2.0.0"
				c.SENDERS.GITHUB.RELEASE.TEMPLATE = "{{.}}"
				c.SENDERS.GITHUB.REPOSITORY.NAME = "test"

				return c
			},
		},
		// Sender stdout
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_FORMAT", "template")
				setenv("CHYLE_SENDERS_STDOUT_TEMPLATE", "{{.}}")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASSTDOUTSENDER = true
				c.FEATURES.HASSENDERS = true
				c.SENDERS.STDOUT.FORMAT = "template"
				c.SENDERS.STDOUT.TEMPLATE = "{{.}}"

				return c
			},
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_STDOUT_FORMAT", "json")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.HASSTDOUTSENDER = true
				c.FEATURES.HASSENDERS = true
				c.SENDERS.STDOUT.FORMAT = "json"

				return c
			},
		},
	}

	for _, test := range tests {
		restoreEnvs()
		chyleConfig = CHYLE{}
		test.f()

		config, err := envh.NewEnvTree("^CHYLE", "_")

		assert.NoError(t, err)

		err = resolveConfig(&config)

		assert.NoError(t, err)

		assert.Equal(t, test.c(), chyleConfig)
	}
}

func TestDebugConfig(t *testing.T) {
	chyleConfig = CHYLE{}
	b := []byte{}

	buffer := bytes.NewBuffer(b)

	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

	EnableDebugging = true

	debugConfig()

	for {
		p := buffer.Next(100)

		if len(p) == 0 {
			break
		}

		b = append(b, p...)
	}

	actual := string(b)

	assert.Regexp(t, `CHYLE - \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} {\n\s+"GIT": {\n\s+"REPOSITORY": {\n`, actual, "Must output given format with argument when debug is enabled")
}

func TestDebugConfigWithDebugDisabled(t *testing.T) {
	chyleConfig = CHYLE{}
	b := []byte{}

	buffer := bytes.NewBuffer(b)

	logger = log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

	EnableDebugging = false

	debugConfig()

	_, err := buffer.ReadString('\n')

	assert.EqualError(t, err, "EOF", "Must return EOF error")
}
