package prompt

import (
	"fmt"
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

const json = "json"
const template = "template"

func newSenders(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		senderChoice,
		newStdoutSender(store),
		newCustomAPISender(store),
		newGithubReleaseSender(store),
	)
}

var senderChoice = []strumt.Prompter{
	builder.NewSwitchPrompt("senderChoice", addMainMenuAndQuitChoice(
		[]builder.SwitchConfig{
			{
				Choice:       "1",
				PromptString: "Add an stdout sender",
				NextPromptID: "senderStdoutFormat",
			},
			{
				Choice:       "2",
				PromptString: "Add a github release sender",
				NextPromptID: "githubReleaseSenderCredentialsToken",
			},
			{
				Choice:       "3",
				PromptString: "Add a custom api sender",
				NextPromptID: "customAPISenderToken",
			},
		},
	)),
}

func newStdoutSender(store *builder.Store) []strumt.Prompter {
	return []strumt.Prompter{
		&builder.GenericPrompt{
			PromptID:  "senderStdoutFormat",
			PromptStr: "Set output format : json or template",
			OnSuccess: func(val string) string {
				if val == json {
					return "senderChoice"
				}
				return "senderStdoutTemplate"
			},
			OnError: func(err error) string {
				return "senderStdoutFormat"
			},
			ParseValue: func(val string) error {
				if val != json && val != template {
					return fmt.Errorf(`"%s" is not a valid format, it must be either "json" or "template"`, val)
				}

				return builder.ParseEnv(func(value string) error { return nil }, "CHYLE_SENDERS_STDOUT_FORMAT", store)(val)
			},
		},
		builder.NewEnvPrompt(
			builder.EnvConfig{
				ID:           "senderStdoutTemplate",
				NextID:       "senderChoice",
				Env:          "CHYLE_SENDERS_STDOUT_TEMPLATE",
				PromptString: "Set a template used to dump to stdout. The syntax follows the golang template (more information here : https://github.com/antham/chyle/wiki/6-Templates)",
				Validator:    validateTemplate,
			}, store,
		),
	}
}

func newGithubReleaseSender(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(githubReleaseSender, store)
}

var githubReleaseSender = []builder.EnvConfig{
	{
		ID:           "githubReleaseSenderCredentialsToken",
		NextID:       "githubReleaseSenderCredentialsOwer",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN",
		PromptString: "Set github oauth token used to publish a release",
		Validator:    validateDefinedValue,
	},
	{
		ID:           "githubReleaseSenderCredentialsOwer",
		NextID:       "githubReleaseSenderRepositoryName",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER",
		PromptString: "Set github owner used in credentials",
		Validator:    validateDefinedValue,
	},
	{
		ID:           "githubReleaseSenderRepositoryName",
		NextID:       "githubReleaseSenderReleaseDraft",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME",
		PromptString: "Set github repository where we will publish the release",
		Validator:    validateDefinedValue,
	},
	{
		ID:           "githubReleaseSenderReleaseDraft",
		NextID:       "githubReleaseSenderReleaseName",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_DRAFT",
		PromptString: "Set if release must be marked as a draft (false or true)",
		Validator:    validateBoolean,
	},
	{
		ID:           "githubReleaseSenderReleaseName",
		NextID:       "githubReleaseSenderReleasePrerelease",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_NAME",
		PromptString: "Set the title of the release, just return if you don't want to give a title",
		Validator:    noOpValidator,
	},
	{
		ID:           "githubReleaseSenderReleasePrerelease",
		NextID:       "githubReleaseSenderReleaseTagName",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_PRERELEASE",
		PromptString: "Set if the release must be marked as a prerelease (false or true)",
		Validator:    validateBoolean,
	},
	{
		ID:           "githubReleaseSenderReleaseTagName",
		NextID:       "githubReleaseSenderReleaseTargetCommit",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME",
		PromptString: "Set release tag to create, when you update a release it will be used to find out the release tied to this tag",
		Validator:    validateDefinedValue,
	},
	{
		ID:           "githubReleaseSenderReleaseTargetCommit",
		NextID:       "githubReleaseSenderReleaseTemplate",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TARGETCOMMITISH",
		PromptString: "Set the commitish value that determines where the git tag must created from",
		Validator:    validateDefinedValue,
	},
	{
		ID:           "githubReleaseSenderReleaseTemplate",
		NextID:       "githubReleaseSenderReleaseUpdate",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE",
		PromptString: "Set a template used to dump the release body. The syntax follows the golang template (more information here : https://github.com/antham/chyle/wiki/6-Templates)",

		Validator: validateTemplate,
	},
	{
		ID:           "githubReleaseSenderReleaseUpdate",
		NextID:       "senderChoice",
		Env:          "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_UPDATE",
		PromptString: "Set to true if you want to update an existing changelog, typical usage would be when you produce a release through GUI github release system",
		Validator:    validateBoolean,
	},
}

func newCustomAPISender(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(customAPISender, store)
}

var customAPISender = []builder.EnvConfig{
	{
		ID:           "customAPISenderToken",
		NextID:       "customAPISenderURL",
		Env:          "CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN",
		PromptString: `Set an access token that would be given in authorization header when calling your API`,
		Validator:    validateDefinedValue,
	},
	{
		ID:           "customAPISenderURL",
		NextID:       "senderChoice",
		Env:          "CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL",
		PromptString: "Set the URL endpoint where the POST request will be sent",
		Validator:    validateURL,
	},
}
