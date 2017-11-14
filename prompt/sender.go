package prompt

import (
	"fmt"
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func newSenders(store *builder.Store) []strumt.Prompter {
	return mergePrompters(
		senderChoice,
		newStdoutSender(store),
		newCustomAPISender(store),
		newGithubReleaseSender(store),
	)
}

var senderChoice = []strumt.Prompter{
	builder.NewSwitchPrompt("senderChoice", addMainMenuAndQuitChoice([]builder.SwitchConfig{{"1", "Add an stdout sender", "senderStdoutFormat"}, {"2", "Add a github release sender", "githubReleaseSenderCredentialsToken"}, {"3", "Add a custom api sender", "customAPISenderToken"}})),
}

func newStdoutSender(store *builder.Store) []strumt.Prompter {
	return []strumt.Prompter{
		builder.NewPromptWithCustomHandlers(
			builder.EnvConfig{"senderStdoutFormat", "", "CHYLE_SENDERS_STDOUT_FORMAT", "Set output format : json or template"},
			func(val string) string {
				if val == "json" {
					return "senderChoice"
				}
				return "senderStdoutTemplate"
			},
			func(err error) string {
				return "senderStdoutFormat"
			},
			func(val string) error {
				if val != "json" && val != "template" {
					return fmt.Errorf(`"%s" is not a valid format, it must be either "json" or "template"`, val)
				}

				return builder.ParseEnv("CHYLE_SENDERS_STDOUT_FORMAT", store)(val)
			},
			store,
		),
		builder.NewEnvPrompt(builder.EnvConfig{"senderStdoutTemplate", "senderChoice", "CHYLE_SENDERS_STDOUT_TEMPLATE", "Set a template following golang template (more information here : https://github.com/antham/chyle#template)"}, store),
	}
}

func newGithubReleaseSender(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(githubReleaseSender, store)
}

var githubReleaseSender = []builder.EnvConfig{
	{"githubReleaseSenderCredentialsToken", "githubReleaseSenderCredentialsOwer", "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OAUTHTOKEN", "Set github oauth token used to publish a release"},
	{"githubReleaseSenderCredentialsOwer", "githubReleaseSenderReleaseDraft", "CHYLE_SENDERS_GITHUBRELEASE_CREDENTIALS_OWNER", "Set github user"},
	{"githubReleaseSenderReleaseDraft", "githubReleaseSenderReleaseName", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_DRAFT", "Set if release must be marked as a draft (false or true)"},
	{"githubReleaseSenderReleaseName", "githubReleaseSenderReleasePrerelease", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_NAME", "Set the title of the release"},
	{"githubReleaseSenderReleasePrerelease", "githubReleaseSenderReleaseTagName", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_PRERELEASE", "Set if the release must be marked as prerelease (false or true)"},
	{"githubReleaseSenderReleaseTagName", "githubReleaseSenderReleaseTargetCommit", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TAGNAME", "Set release tag to create, when you update a release it will be used to find out release tied to this tag"},
	{"githubReleaseSenderReleaseTargetCommit", "githubReleaseSenderReleaseTemplate", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TARGETCOMMITISH", "Set the commitish value that determines where the git tag is created from"},
	{"githubReleaseSenderReleaseTemplate", "githubReleaseSenderReleaseUpdate", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_TEMPLATE", "Set a template following golang template (more information here : https://github.com/antham/chyle#template)"},
	{"githubReleaseSenderReleaseUpdate", "githubReleaseSenderRepositoryName", "CHYLE_SENDERS_GITHUBRELEASE_RELEASE_UPDATE", "Set to true if you want to update an existing changelog, typical usage would be when you produce a release through GUI github release system"},
	{"githubReleaseSenderRepositoryName", "senderChoice", "CHYLE_SENDERS_GITHUBRELEASE_REPOSITORY_NAME", "Set github repository where we will publish the release (false or true)"},
}

func newCustomAPISender(store *builder.Store) []strumt.Prompter {
	return builder.NewEnvPrompts(customAPISender, store)
}

var customAPISender = []builder.EnvConfig{
	{"customAPISenderToken", "customAPISenderURL", "CHYLE_SENDERS_CUSTOMAPI_CREDENTIALS_TOKEN", `Set an access token that would be given in request header "Authorization" to API`},
	{"customAPISenderURL", "senderChoice", "CHYLE_SENDERS_CUSTOMAPI_ENDPOINT_URL", "Set the URL endpoint where the POST request will be sent"},
}
