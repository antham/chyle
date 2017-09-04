package prompt

import (
	"github.com/antham/strumt"
)

func newDecorators(store *Store) []strumt.Prompter {
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
	&switchPrompt{
		"decoratorChoice",
		addMainMenuAndQuitChoice(
			[]switchChoice{
				switchChoice{"1", "Add a custom api decorator", "extractorOrigKeyCustomAPI"},
				switchChoice{"2", "Add a jira issue decorator", "extractorOrigKeyJiraIssueID"},
				switchChoice{"3", "Add a github issue decorator", "extractorOrigKeyGithubIssueID"},
				switchChoice{"4", "Add a shell decorator", "shellDecoratorCommand"},
				switchChoice{"5", "Add an environment variable decorator", "envDecoratorVarName"},
			},
		),
	},
}

func newCustomAPIDecorator(store *Store) []strumt.Prompter {
	return mergePrompters(
		newGroupEnvPromptWithCounter(customAPIDecoratorKeys, store),
		newEnvPrompts(customAPIDecorator, store),
		customAPIDecoratorChoice,
	)
}

var customAPIDecorator = []envConfig{
	envConfig{"extractorOrigKeyCustomAPI", "extractorDestKeyCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_ORIGKEY", "Enter a commit field where your custom api id is located"},
	envConfig{"extractorDestKeyCustomAPI", "extractorRegCustomAPI", "CHYLE_EXTRACTORS_CUSTOMAPIID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	envConfig{"extractorRegCustomAPI", "customAPIDecoratorEndpoint", "CHYLE_EXTRACTORS_CUSTOMAPIID_REG", "Enter a regexp to extract custom api id"},
	envConfig{"customAPIDecoratorEndpoint", "customAPIDecoratorToken", "CHYLE_DECORATORS_CUSTOMAPIID_ENDPOINT_URL", "Enter custom api endpoint URL, use {{ID}} as a placeholder to interpolate the id you extracted before in URL if you need to"},
	envConfig{"customAPIDecoratorToken", "customAPIDecoratorDestKey", "CHYLE_DECORATORS_CUSTOMAPIID_CREDENTIALS_TOKEN", "Enter token submitted as authorization header when calling your api"},
}

var customAPIDecoratorKeys = []envConfig{
	envConfig{"customAPIDecoratorDestKey", "customAPIDecoratorField", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_DESTKEY", "A name for the key which will receive the data extracted from the custom api"},
	envConfig{"customAPIDecoratorField", "customAPIDecoratorChoice", "CHYLE_DECORATORS_CUSTOMAPIID_KEYS_*_FIELD", `The field to extract from your custom api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var customAPIDecoratorChoice = []strumt.Prompter{
	&switchPrompt{
		"customAPIDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]switchChoice{
				switchChoice{
					"1", "Add a new custom api decorator field", "customAPIDecoratorDestKey",
				},
			},
		),
	},
}

func newJiraIssueDecorator(store *Store) []strumt.Prompter {
	return mergePrompters(
		newGroupEnvPromptWithCounter(jiraIssueDecoratorKeys, store),
		newEnvPrompts(jiraIssueDecorator, store),
		jiraIssueDecoratorChoice,
	)
}

var jiraIssueDecorator = []envConfig{
	envConfig{"extractorOrigKeyJiraIssueID", "extractorDestKeyJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_ORIGKEY", "Enter a commit field where your jira issue id is located"},
	envConfig{"extractorDestKeyJiraIssueID", "extractorRegJiraIssueID", "CHYLE_EXTRACTORS_JIRAISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	envConfig{"extractorRegJiraIssueID", "jiraIssueDecoratorEndpoint", "CHYLE_EXTRACTORS_JIRAISSUEID_REG", "Enter a regexp to extract jira issue id"},
	envConfig{"jiraIssueDecoratorEndpoint", "jiraIssueDecoratorUsername", "CHYLE_DECORATORS_JIRAISSUE_ENDPOINT_URL", "Enter jira api endpoint URL"},
	envConfig{"jiraIssueDecoratorUsername", "jiraIssueDecoratorPassword", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_USERNAME", "Enter jira username"},
	envConfig{"jiraIssueDecoratorPassword", "jiraIssueDecoratorDestKey", "CHYLE_DECORATORS_JIRAISSUE_CREDENTIALS_PASSWORD", "Enter jira password"},
}

var jiraIssueDecoratorKeys = []envConfig{
	envConfig{"jiraIssueDecoratorDestKey", "jiraIssueDecoratorField", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	envConfig{"jiraIssueDecoratorField", "jiraIssueDecoratorChoice", "CHYLE_DECORATORS_JIRAISSUE_KEYS_*_FIELD", `The field to extract from jira api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var jiraIssueDecoratorChoice = []strumt.Prompter{
	&switchPrompt{
		"jiraIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]switchChoice{
				switchChoice{
					"1", "Add a new jira issue decorator field", "jiraIssueDecoratorDestKey",
				},
			},
		),
	},
}

func newGithubIssueDecorator(store *Store) []strumt.Prompter {
	return mergePrompters(
		newGroupEnvPromptWithCounter(githubIssueDecoratorKeys, store),
		newEnvPrompts(githubIssueDecorator, store),
		githubIssueDecoratorChoice,
	)
}

var githubIssueDecorator = []envConfig{
	envConfig{"extractorOrigKeyGithubIssueID", "extractorDestKeyGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_ORIGKEY", "Enter a commit field where your github issue id is located"},
	envConfig{"extractorDestKeyGithubIssueID", "extractorRegGithubIssueID", "CHYLE_EXTRACTORS_GITHUBISSUEID_DESTKEY", "Enter a name for the key which will receive the extracted value"},
	envConfig{"extractorRegGithubIssueID", "githubIssueCredentialsToken", "CHYLE_EXTRACTORS_GITHUBISSUEID_REG", "Enter a regexp to extract github issue id"},
	envConfig{"githubIssueCredentialsToken", "githubIssueCredentialsOwner", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OAUTHTOKEN", "Enter github oauth token"},
	envConfig{"githubIssueCredentialsOwner", "githubIssueDecoratorDestKey", "CHYLE_DECORATORS_GITHUBISSUE_CREDENTIALS_OWNER", "Enter github owner name"},
}

var githubIssueDecoratorKeys = []envConfig{
	envConfig{"githubIssueDecoratorDestKey", "githubIssueDecoratorField", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_DESTKEY", "A name for the key which will receive the extracted value"},
	envConfig{"githubIssueDecoratorField", "githubIssueDecoratorChoice", "CHYLE_DECORATORS_GITHUBISSUE_KEYS_*_FIELD", `The field to extract from github issue api response payload, use dot notation to extract a deep value (eg: "fields.summary")`},
}

var githubIssueDecoratorChoice = []strumt.Prompter{
	&switchPrompt{
		"githubIssueDecoratorChoice",
		addMainMenuAndQuitChoice(
			[]switchChoice{
				switchChoice{
					"1", "Add a new github issue decorator field", "githubIssueDecoratorDestKey",
				},
			},
		),
	},
}

func newShellDecorator(store *Store) []strumt.Prompter {
	return newGroupEnvPromptWithCounter(shellDecoratorKeys, store)
}

var shellDecoratorKeys = []envConfig{
	envConfig{"shellDecoratorCommand", "shellDecoratorOrigKey", "CHYLE_DECORATORS_SHELL_*_COMMAND", "Shell command to execute"},
	envConfig{"shellDecoratorOrigKey", "shellDecoratorDestKey", "CHYLE_DECORATORS_SHELL_*_ORIGKEY", "A field from which you want to use the content to pipe a command on"},
	envConfig{"shellDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_SHELL_*_DESTKEY", "A name for the key which will receive the extracted value"},
}

func newEnvDecorator(store *Store) []strumt.Prompter {
	return newGroupEnvPromptWithCounter(envDecoratorKeys, store)
}

var envDecoratorKeys = []envConfig{
	envConfig{"envDecoratorVarName", "envDecoratorDestKey", "CHYLE_DECORATORS_ENV_*_VARNAME", "Environment variable name to dump in metadatas"},
	envConfig{"envDecoratorDestKey", "decoratorChoice", "CHYLE_DECORATORS_ENV_*_DESTKEY", "The name of the key where to store dumped value in metadatas"},
}
