package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/antham/envh"
	"github.com/stretchr/testify/assert"

	"github.com/antham/chyle/chyle/matchers"
)

var envs map[string]string

func TestMain(m *testing.M) {
	saveExistingEnvs()
	code := m.Run()
	os.Exit(code)
}

func saveExistingEnvs() {
	var err error
	env := envh.NewEnv()

	envs, err = env.FindEntries(".*")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func restoreEnvs() {
	os.Clearenv()

	if len(envs) != 0 {
		for key, value := range envs {
			setenv(key, value)
		}
	}
}

func setenv(key string, value string) {
	err := os.Setenv(key, value)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestCreateWithErrors(t *testing.T) {
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
		// Decorators custom api
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_TEST", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test")
			},
			`environment variable missing : "CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
			},
			`environment variable missing : "CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
			},
			`provide a valid URL for "CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "test" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
			},
			`define at least one environment variable couple "CHYLE_DECORATORS_CUSTOMAPI_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_CUSTOMAPI_KEYS_*_FIELD", replace "*" with your own naming`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "CHYLE_EXTRACTORS_CUSTOMAPIID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "CHYLE_EXTRACTORS_CUSTOMAPIID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com/get")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "customApiId")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_CUSTOMAPIID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com/get/{{ID}}")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_REG", "test")
			},
			`variable CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY must be equal to "customApiId"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com/get")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "customApiId")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_REG", "test")
			},
			`ensure you defined a placeholder {{ID}} in URL defined in "CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL"`,
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
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
			},
			`environments variables missing : "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "testcom")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
			},
			`provide a valid URL for "CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "testcom" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
			},
			`define at least one environment variable couple "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_FIELD", replace "*" with your own naming`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "jiraIssueId")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "jiraIssueId")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "jiraIssueId")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_JIRAISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_REG", "test")
			},
			`variable CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY must be equal to "jiraIssueId"`,
		},
		// Decorator github
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environments variables missing : "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
			},
			`define at least one environment variable couple "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_FIELD", replace "*" with your own naming`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "githubIssueId")
			},
			`environment variable missing : "CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "githubIssueId")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "githubIssueId")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "githubIssueId")
			},
			`environment variable missing : "CHYLE_EXTRACTORS_GITHUBISSUEID_REG"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_REG", "test")
			},
			`variable CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY must be equal to "githubIssueId"`,
		},
		// Decorator shell
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_COMMAND", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "CHYLE_DECORATORS_SHELL_TEST_ORIGKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_ORIGKEY", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "CHYLE_DECORATORS_SHELL_TEST_COMMAND"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "test")
			},
			`environments variables missing : "CHYLE_DECORATORS_SHELL_TEST_ORIGKEY", "CHYLE_DECORATORS_SHELL_TEST_COMMAND"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_COMMAND", "test")
				setenv("CHYLE_DECORATORS_SHELL_TEST_ORIGKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_SHELL_TEST_DESTKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_COMMAND", "test")
				setenv("CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_SHELL_TEST_ORIGKEY"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_ORIGKEY", "test")
				setenv("CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "test")
			},
			`environment variable missing : "CHYLE_DECORATORS_SHELL_TEST_COMMAND"`,
		},
		// Sender github
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_TEST", "test")
			},
			`environments variables missing : "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environments variables missing : "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "v2.0.0")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE", "{{.....}}")
			},
			`provide a valid template string for "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE" : "template: test:1: unexpected <.> in operand", "{{.....}}" given`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE", "{{.}}")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE", "{{.}}")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "v2.0.0")
			},
			`environment variable missing : "CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME"`,
		},
		// Sender custom api
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TEST", "test")
			},
			`environment variable missing : "CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			`environment variable missing : "CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL", "http://test.com")
			},
			`environment variable missing : "CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN"`,
		},
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL", "test")
			},
			`provide a valid URL for "CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL", "test" given`,
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

		_, err = Create(&config)

		errDetail := fmt.Sprintf("Test %d failed", i+1)

		assert.Error(t, err, errDetail)
		assert.EqualError(t, err, test.e, errDetail)
	}
}

func TestCreate(t *testing.T) {
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
				c.FEATURES.MATCHERS.ENABLED = true
				c.FEATURES.MATCHERS.AUTHOR = true
				c.FEATURES.MATCHERS.COMMITTER = true
				c.FEATURES.MATCHERS.TYPE = true
				c.FEATURES.MATCHERS.MESSAGE = true
				c.MATCHERS = matchers.Config{
					MESSAGE:   regexp.MustCompile(".*"),
					AUTHOR:    regexp.MustCompile(".*"),
					COMMITTER: regexp.MustCompile(".*"),
					TYPE:      "regular",
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
				c.FEATURES.EXTRACTORS.ENABLED = true
				c.EXTRACTORS = map[string]struct {
					ORIGKEY string
					DESTKEY string
					REG     *regexp.Regexp
				}{
					"TEST": {
						"test",
						"test",
						regexp.MustCompile(".*"),
					},
				}

				return c
			},
		},
		// Decorators custom api
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_ENDPOINT_URL", "http://test.com/get/{{ID}}")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_CREDENTIALS_TOKEN", "da39a3ee5e6b4b0d3255bfef95601890afd80709")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_DESTKEY", "destKey")
				setenv("CHYLE_DECORATORS_CUSTOMAPI_KEYS_TEST_FIELD", "field")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "customApiId")
				setenv("CHYLE_EXTRACTORS_CUSTOMAPIID_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.EXTRACTORS.ENABLED = true
				c.FEATURES.DECORATORS.ENABLED = true
				c.FEATURES.DECORATORS.CUSTOMAPI = true
				c.EXTRACTORS = map[string]struct {
					ORIGKEY string
					DESTKEY string
					REG     *regexp.Regexp
				}{
					"CUSTOMAPIID": {
						"test",
						"customApiId",
						regexp.MustCompile(".*"),
					},
				}
				c.DECORATORS.CUSTOMAPI.CREDENTIALS.TOKEN = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
				c.DECORATORS.CUSTOMAPI.ENDPOINT.URL = "http://test.com/get/{{ID}}"
				c.DECORATORS.CUSTOMAPI.KEYS = map[string]struct {
					DESTKEY string
					FIELD   string
				}{
					"TEST": {
						"destKey",
						"field",
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
				c.FEATURES.DECORATORS.ENABLED = true
				c.FEATURES.DECORATORS.ENV = true
				c.DECORATORS.ENV = map[string]struct {
					DESTKEY string
					VARNAME string
				}{
					"TEST": {
						"destkey",
						"var",
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
				setenv("CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "http://test.com")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "test")
				setenv("CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "password")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_DESTKEY", "destkey")
				setenv("CHYLE_DECORATORS_JIRAISSUE_KEYS_TEST_FIELD", "field")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "jiraIssueId")
				setenv("CHYLE_EXTRACTORS_JIRAISSUEID_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.EXTRACTORS.ENABLED = true
				c.FEATURES.DECORATORS.ENABLED = true
				c.FEATURES.DECORATORS.JIRAISSUE = true
				c.EXTRACTORS = map[string]struct {
					ORIGKEY string
					DESTKEY string
					REG     *regexp.Regexp
				}{
					"JIRAISSUEID": {
						"test",
						"jiraIssueId",
						regexp.MustCompile(".*"),
					},
				}
				c.DECORATORS.JIRAISSUE.ENDPOINT.URL = "http://test.com"
				c.DECORATORS.JIRAISSUE.CREDENTIALS.USERNAME = "test"
				c.DECORATORS.JIRAISSUE.CREDENTIALS.PASSWORD = "password"
				c.DECORATORS.JIRAISSUE.KEYS = map[string]struct {
					DESTKEY string
					FIELD   string
				}{
					"TEST": {
						"destkey",
						"field",
					},
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
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_REPOSITORY_NAME", "test")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_DESTKEY", "destkey")
				setenv("CHYLE_DECORATORS_GITHUBISSUE_KEYS_TEST_FIELD", "field")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "test")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "githubIssueId")
				setenv("CHYLE_EXTRACTORS_GITHUBISSUEID_REG", ".*")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.EXTRACTORS.ENABLED = true
				c.FEATURES.DECORATORS.ENABLED = true
				c.FEATURES.DECORATORS.GITHUBISSUE = true
				c.EXTRACTORS = map[string]struct {
					ORIGKEY string
					DESTKEY string
					REG     *regexp.Regexp
				}{
					"GITHUBISSUEID": {
						"test",
						"githubIssueId",
						regexp.MustCompile(".*"),
					},
				}
				c.DECORATORS.GITHUBISSUE.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
				c.DECORATORS.GITHUBISSUE.CREDENTIALS.OWNER = "test"
				c.DECORATORS.GITHUBISSUE.REPOSITORY.NAME = "test"
				c.DECORATORS.GITHUBISSUE.KEYS = map[string]struct {
					DESTKEY string
					FIELD   string
				}{
					"TEST": {
						"destkey",
						"field",
					},
				}

				return c
			},
		},
		// Decorator shell
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_DECORATORS_SHELL_TEST_COMMAND", "test")
				setenv("CHYLE_DECORATORS_SHELL_TEST_ORIGKEY", "test")
				setenv("CHYLE_DECORATORS_SHELL_TEST_DESTKEY", "test")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.DECORATORS.ENABLED = true
				c.FEATURES.DECORATORS.SHELL = true
				c.DECORATORS.SHELL = map[string]struct {
					COMMAND string
					ORIGKEY string
					DESTKEY string
				}{
					"TEST": {
						"test",
						"test",
						"test",
					},
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
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "user")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "d41d8cd98f00b204e9800998ecf8427e")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE", "{{.}}")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "v2.0.0")
				setenv("CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME", "test")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.SENDERS.ENABLED = true
				c.FEATURES.SENDERS.GITHUBRELEASE = true
				c.SENDERS.GITHUBRELEASE.CREDENTIALS.OAUTHTOKEN = "d41d8cd98f00b204e9800998ecf8427e"
				c.SENDERS.GITHUBRELEASE.CREDENTIALS.OWNER = "user"
				c.SENDERS.GITHUBRELEASE.RELEASE.TAGNAME = "v2.0.0"
				c.SENDERS.GITHUBRELEASE.RELEASE.TEMPLATE = "{{.}}"
				c.SENDERS.GITHUBRELEASE.REPOSITORY.NAME = "test"

				return c
			},
		},
		// Sender custom api
		{
			func() {
				setenv("CHYLE_GIT_REPOSITORY_PATH", "test")
				setenv("CHYLE_GIT_REFERENCE_FROM", "v1.0.0")
				setenv("CHYLE_GIT_REFERENCE_TO", "v2.0.0")
				setenv("CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL", "http://test.com/releases")
				setenv("CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN", "d41d8cd98f00b204e9800998ecf8427e")
			},
			func() CHYLE {
				c := CHYLE{}
				c.GIT.REPOSITORY.PATH = "test"
				c.GIT.REFERENCE.FROM = "v1.0.0"
				c.GIT.REFERENCE.TO = "v2.0.0"
				c.FEATURES.SENDERS.ENABLED = true
				c.FEATURES.SENDERS.CUSTOMAPI = true
				c.SENDERS.CUSTOMAPI.CREDENTIALS.TOKEN = "d41d8cd98f00b204e9800998ecf8427e"
				c.SENDERS.CUSTOMAPI.ENDPOINT.URL = "http://test.com/releases"

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
				c.FEATURES.SENDERS.STDOUT = true
				c.FEATURES.SENDERS.ENABLED = true
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
				c.FEATURES.SENDERS.STDOUT = true
				c.FEATURES.SENDERS.ENABLED = true
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

		actual, err := Create(&config)

		assert.NoError(t, err)

		assert.Equal(t, test.c(), *actual)
	}
}

func TestDebug(t *testing.T) {
	chyleConfig = CHYLE{}
	b := []byte{}

	buffer := bytes.NewBuffer(b)

	logger := log.New(buffer, "CHYLE - ", log.Ldate|log.Ltime)

	Debug(&chyleConfig, logger)

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
