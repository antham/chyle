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
				builder.SwitchConfig{"1", "Add a custom api decorator", "extractorOrigKeyCustomAPI"},
				builder.SwitchConfig{"2", "Add a jira issue decorator", "extractorOrigKeyJiraIssueID"},
				builder.SwitchConfig{"3", "Add a github issue decorator", "extractorOrigKeyGithubIssueID"},
				builder.SwitchConfig{"4", "Add a shell decorator", "shellDecoratorCommand"},
				builder.SwitchConfig{"5", "Add an environment variable decorator", "envDecoratorVarName"},
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
	builder.EnvConfig{"extractorOrigKeyCustomAPI", "extractorDestKeyCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "Enter a commit field where your custom api id is located"},
	builder.EnvConfig{"extractorDestKeyCustomAPI", "extractorRegCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	builder.EnvConfig{"extractorRegCustomAPI", "customAPIDecoratorEndpoint", "CHYLE_EXTRACTORS_CUSTOMAPIID_REG", "Enter a regexp to extract custom api id"},
	builder.EnvConfig{"customAPIDecoratorEndpoint", "customAPIDecoratorToken", "CHYLE_DECORATORS_CUSTOMAPIID_ENDPOINT_URL", "Enter custom api endpoint URL, use {{ID}} as a placeholder to interpolate the id you extracted before in URL if you need to"},
	builder.EnvConfig{"customAPIDecoratorToken", "customAPIDecoratorDestKey", "CHYLE_DECORATORS_CUSTOMAPIID_CREDENTIALS_TOKEN", "Enter token submitted as authorization header when calling your api"},
}

var customAPIDecoratorKeys = []builder.EnvConfig{
	builder.EnvConfig{"customAPIDecoratorDestKey", "customAPIDecoratorField", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_DESTKEY", "A name for the key which will receive the data extracted from the custom api"},
	builder.EnvConfig{"customAPIDecoratorField", "customAPIDecoratorChoice", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_FIELD", `The field to extract from your custom api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var customAPIDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"customAPIDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				builder.SwitchConfig{
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
	builder.EnvConfig{"extractorOrigKeyJiraIssueID", "extractorDestKeyJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "Enter a commit field where your jira issue id is located"},
	builder.EnvConfig{"extractorDestKeyJiraIssueID", "extractorRegJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	builder.EnvConfig{"extractorRegJiraIssueID", "jiraIssueDecoratorEndpoint", "CHYLE_EXTRACTORS_JIRAISSUEID_REG", "Enter a regexp to extract jira issue id"},
	builder.EnvConfig{"jiraIssueDecoratorEndpoint", "jiraIssueDecoratorUsername", "CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "Enter jira api endpoint URL"},
	builder.EnvConfig{"jiraIssueDecoratorUsername", "jiraIssueDecoratorPassword", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "Enter jira username"},
	builder.EnvConfig{"jiraIssueDecoratorPassword", "jiraIssueDecoratorDestKey", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "Enter jira password"},
}

var jiraIssueDecoratorKeys = []builder.EnvConfig{
	builder.EnvConfig{"jiraIssueDecoratorDestKey", "jiraIssueDecoratorField", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	builder.EnvConfig{"jiraIssueDecoratorField", "jiraIssueDecoratorChoice", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_FIELD", `The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var jiraIssueDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"jiraIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				builder.SwitchConfig{
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
	builder.EnvConfig{"extractorOrigKeyGithubIssueID", "extractorDestKeyGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "Enter a commit field where your github issue id is located"},
	builder.EnvConfig{"extractorDestKeyGithubIssueID", "extractorRegGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	builder.EnvConfig{"extractorRegGithubIssueID", "githubIssueCredentialsToken", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG", "Enter a regexp to extract github issue id"},
	builder.EnvConfig{"githubIssueCredentialsToken", "githubIssueCredentialsOwner", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "Enter github oauth token"},
	builder.EnvConfig{"githubIssueCredentialsOwner", "githubIssueDecoratorDestKey", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "Enter github owner name"},
}

var githubIssueDecoratorKeys = []builder.EnvConfig{
	builder.EnvConfig{"githubIssueDecoratorDestKey", "githubIssueDecoratorField", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	builder.EnvConfig{"githubIssueDecoratorField", "githubIssueDecoratorChoice", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_FIELD", `The field to extract from github issue api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var githubIssueDecoratorChoice = []strumt.Prompter{
	builder.NewSwitchPrompt(
		"githubIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]builder.SwitchConfig{
				builder.SwitchConfig{
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
	builder.EnvConfig{"shellDecoratorCommand", "shellDecoratorOrigKey", "CHYLE_DECORATORS_SHELL_*_COMMAND", "Shell command to execute"},
	builder.EnvConfig{"shellDecoratorOrigKey", "shellDecoratorDestKey", "CHYLE_DECORATORS_SHELL_*_ORIGKEY", "A field from which you want to use the content to pipe a command on"},
	builder.EnvConfig{"shellDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_SHELL_*_DESTKEY", "A name for the key which will receive the extracted value"},
}

func newEnvDecorator(store *builder.Store) []strumt.Prompter {
	return builder.NewGroupEnvPromptWithCounter(envDecoratorKeys, store)
}

var envDecoratorKeys = []builder.EnvConfig{
	builder.EnvConfig{"envDecoratorVarName", "envDecoratorDestKey", "CHYLE_DECORATORS_ENV_*_VARNAME", "Environment variable name to dump in metadatas"},
	builder.EnvConfig{"envDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_ENV_*_DESTKEY", "The name of the key where to store dumped value in metadatas"},
}
