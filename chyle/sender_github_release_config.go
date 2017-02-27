package chyle

import (
	"net/http"

	"github.com/antham/envh"
)

// githubRelease follows https://developer.github.com/v3/repos/releases/#create-a-release
type githubRelease struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft,omitempty"`
	PreRelease      bool   `json:"prerelease,omitempty"`
}

// githubReleaseConfig stores config metadata (github token, request type, and so on...)
type githubReleaseConfig struct {
	oauthToken     string
	owner          string
	repositoryName string
	template       string
	update         bool
}

// buildGithubReleaseSender create a new GithubReleaseSender object from viper config
func buildGithubReleaseSender(config *envh.EnvTree) (sender, error) {
	grConfig, err := buildGithubReleaseConfig(config)

	if err != nil {
		return githubReleaseSender{}, err
	}

	gRelease, err := buildGithubRelease(config)

	if err != nil {
		return githubReleaseSender{}, err
	}

	return newGithubReleaseSender(&http.Client{}, grConfig, gRelease), nil
}

func buildGithubReleaseConfig(config *envh.EnvTree) (githubReleaseConfig, error) {
	grConfig := githubReleaseConfig{}

	err := extractStringConfig(
		config,
		[]strConfigMapping{
			strConfigMapping{
				[]string{"CREDENTIALS", "OAUTHTOKEN"},
				&grConfig.oauthToken,
				true,
			},
			strConfigMapping{
				[]string{"CREDENTIALS", "OWNER"},
				&grConfig.owner,
				true,
			},
			strConfigMapping{
				[]string{"REPOSITORY", "NAME"},
				&grConfig.repositoryName,
				true,
			},
			strConfigMapping{
				[]string{"RELEASE", "TEMPLATE"},
				&grConfig.template,
				true,
			},
		},
		[]string{"SENDERS", "GITHUB"},
	)

	if err != nil {
		return grConfig, err
	}

	err = extractBoolConfig(
		config,
		[]boolConfigMapping{
			boolConfigMapping{
				[]string{"RELEASE", "UPDATE"},
				&grConfig.update,
				false,
			},
		},
		[]string{"SENDERS", "GITHUB"},
	)

	return grConfig, err
}

func buildGithubRelease(config *envh.EnvTree) (githubRelease, error) {
	gRelease := githubRelease{}

	err := extractStringConfig(
		config,
		[]strConfigMapping{
			strConfigMapping{
				[]string{"RELEASE", "TAGNAME"},
				&gRelease.TagName,
				true,
			},
			strConfigMapping{
				[]string{"RELEASE", "TARGETCOMMITISH"},
				&gRelease.TargetCommitish,
				false,
			},
			strConfigMapping{
				[]string{"RELEASE", "NAME"},
				&gRelease.Name,
				false,
			},
		},
		[]string{"SENDERS", "GITHUB"},
	)

	if err != nil {
		return gRelease, err
	}

	err = extractBoolConfig(
		config,
		[]boolConfigMapping{
			boolConfigMapping{
				[]string{"RELEASE", "DRAFT"},
				&gRelease.Draft,
				false,
			},
			boolConfigMapping{
				[]string{"RELEASE", "PRERELEASE"},
				&gRelease.PreRelease,
				false,
			},
		},
		[]string{"SENDERS", "GITHUB"},
	)

	return gRelease, err
}
