package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newDecorators(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		decorator,
		newCustomAPIDecorator(store),
		newJiraIssueDecorator(store),
		newGithubIssueDecorator(store),
		newShellDecorator(store),
		newEnvDecorator(store),
	)
}

var decorator = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"decoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				{"1", "Add a custom api decorator", "extractorOrigKeyCustomAPI"},
				{"2", "Add a jira issue decorator", "extractorOrigKeyJiraIssueID"},
				{"3", "Add a github issue decorator", "extractorOrigKeyGithubIssueID"},
				{"4", "Add a shell decorator", "shellDecoratorCommand"},
				{"5", "Add an environment variable decorator", "envDecoratorVarName"},
			},
		),
	),
}

func newCustomAPIDecorator(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		builder.NewGroupEnvPromptWithCounter(customAPIDecoratorKeys, store),
		builder.NewEnvPrompts(customAPIDecorator, store),
		customAPIDecoratorChoice,
	)
}

var customAPIDecorator = []builder.EnvConfig{
	{"extractorOrigKeyCustomAPI", "extractorDestKeyCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "Enter a commit field where your custom api id is located"},
	{"extractorDestKeyCustomAPI", "extractorRegCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	{"extractorRegCustomAPI", "customAPIDecoratorEndpoint", "CHYLE_EXTRACTORS_CUSTOMAPIID_REG", "Enter a regexp to extract custom api id"},
	{"customAPIDecoratorEndpoint", "customAPIDecoratorToken", "CHYLE_DECORATORS_CUSTOMAPIID_ENDPOINT_URL", "Enter custom api endpoint URL, use {{ID}} as a placeholder to interpolate the id you extracted before in URL if you need to"},
	{"customAPIDecoratorToken", "customAPIDecoratorDestKey", "CHYLE_DECORATORS_CUSTOMAPIID_CREDENTIALS_TOKEN", "Enter token submitted as authorization header when calling your api"},
}

var customAPIDecoratorKeys = []builder.EnvConfig{
	{"customAPIDecoratorDestKey", "customAPIDecoratorField", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_DESTKEY", "A name for the key which will receive the data extracted from the custom api"},
	{"customAPIDecoratorField", "customAPIDecoratorChoice", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_FIELD", `The field to extract from your custom api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var customAPIDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"customAPIDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				{
					"1", "Add a new custom api decorator field", "customAPIDecoratorDestKey",
				},
			},
		),
	),
}

func newJiraIssueDecorator(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		builder.NewGroupEnvPromptWithCounter(jiraIssueDecoratorKeys, store),
		builder.NewEnvPrompts(jiraIssueDecorator, store),
		jiraIssueDecoratorChoice,
	)
}

var jiraIssueDecorator = []builder.EnvConfig{
	{"extractorOrigKeyJiraIssueID", "extractorDestKeyJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "Enter a commit field where your jira issue id is located"},
	{"extractorDestKeyJiraIssueID", "extractorRegJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	{"extractorRegJiraIssueID", "jiraIssueDecoratorEndpoint", "CHYLE_EXTRACTORS_JIRAISSUEID_REG", "Enter a regexp to extract jira issue id"},
	{"jiraIssueDecoratorEndpoint", "jiraIssueDecoratorUsername", "CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "Enter jira api endpoint URL"},
	{"jiraIssueDecoratorUsername", "jiraIssueDecoratorPassword", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "Enter jira username"},
	{"jiraIssueDecoratorPassword", "jiraIssueDecoratorDestKey", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "Enter jira password"},
}

var jiraIssueDecoratorKeys = []builder.EnvConfig{
	{"jiraIssueDecoratorDestKey", "jiraIssueDecoratorField", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	{"jiraIssueDecoratorField", "jiraIssueDecoratorChoice", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_FIELD", `The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var jiraIssueDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"jiraIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				{
					"1", "Add a new jira issue decorator field", "jiraIssueDecoratorDestKey",
				},
			},
		),
	),
}

func newGithubIssueDecorator(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		builder.NewGroupEnvPromptWithCounter(githubIssueDecoratorKeys, store),
		builder.NewEnvPrompts(githubIssueDecorator, store),
		githubIssueDecoratorChoice,
	)
}

var githubIssueDecorator = []builder.EnvConfig{
	{"extractorOrigKeyGithubIssueID", "extractorDestKeyGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "Enter a commit field where your github issue id is located"},
	{"extractorDestKeyGithubIssueID", "extractorRegGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	{"extractorRegGithubIssueID", "githubIssueCredentialsToken", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG", "Enter a regexp to extract github issue id"},
	{"githubIssueCredentialsToken", "githubIssueCredentialsOwner", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "Enter github oauth token"},
	{"githubIssueCredentialsOwner", "githubIssueDecoratorDestKey", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "Enter github owner name"},
}

var githubIssueDecoratorKeys = []builder.EnvConfig{
	{"githubIssueDecoratorDestKey", "githubIssueDecoratorField", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	{"githubIssueDecoratorField", "githubIssueDecoratorChoice", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_FIELD", `The field to extract from github issue api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var githubIssueDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"githubIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				{
					"1", "Add a new github issue decorator field", "githubIssueDecoratorDestKey",
				},
			},
		),
	),
}

func newShellDecorator(store *builder.Store) []strumt.Prompter {
	return builder.NewGroupEnvPromptWithCounter(shellDecoratorKeys, store)
}

var shellDecoratorKeys = []builder.EnvConfig{
	{"shellDecoratorCommand", "shellDecoratorOrigKey", "CHYLE_DECORATORS_SHELL_*_COMMAND", "Shell command to execute"},
	{"shellDecoratorOrigKey", "shellDecoratorDestKey", "CHYLE_DECORATORS_SHELL_*_ORIGKEY", "A field from which you want to use the content to pipe a command on"},
	{"shellDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_SHELL_*_DESTKEY", "A name for the key which will receive the extracted value"},
}

func newEnvDecorator(store *builder.Store) []strumt.Prompter {
	return builder.NewGroupEnvPromptWithCounter(envDecoratorKeys, store)
}

var envDecoratorKeys = []builder.EnvConfig{
	{"envDecoratorVarName", "envDecoratorDestKey", "CHYLE_DECORATORS_ENV_*_VARNAME", "Environment variable name to dump in metadatas"},
	{"envDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_ENV_*_DESTKEY", "The name of the key where to store dumped value in metadatas"},
}
