package chyle

import (
	"encoding/json"
	"strings"

	"github.com/antham/envh"
)

// configurater must be implemented to process custom config
type configurater interface {
	process(config *CHYLE) (bool, error)
}

var chyleConfig CHYLE

// codebeat:disable[TOO_MANY_IVARS]

// CHYLE hold config extracted from environment variables
type CHYLE struct {
	FEATURES struct {
		HASMATCHERS            bool
		HASEXTRACTORS          bool
		HASDECORATORS          bool
		HASSENDERS             bool
		HASJIRADECORATOR       bool
		HASENVDECORATOR        bool
		HASGITHUBRELEASESENDER bool
		HASSTDOUTSENDER        bool
	} `json:"-"`
	GIT struct {
		REPOSITORY struct {
			PATH string
		}
		REFERENCE struct {
			FROM string
			TO   string
		}
	}
	MATCHERS   map[string]string
	EXTRACTORS map[string]map[string]string
	DECORATORS struct {
		JIRA struct {
			CREDENTIALS struct {
				URL      string
				USERNAME string
				PASSWORD string
			}
			KEYS map[string]string
		}
		ENV map[string]map[string]string
	}
	SENDERS struct {
		STDOUT struct {
			FORMAT   string
			TEMPLATE string
		}
		GITHUB struct {
			CREDENTIALS struct {
				OAUTHTOKEN string
				OWNER      string
			}
			RELEASE struct {
				DRAFT           bool
				NAME            string
				PRERELEASE      bool
				TAGNAME         string
				TARGETCOMMITISH string
				TEMPLATE        string
				UPDATE          bool
			}
			REPOSITORY struct {
				NAME string
			}
		}
	}
}

// codebeat:enable[TOO_MANY_IVARS]

// Walk traverses struct to populate or validate fields
func (c *CHYLE) Walk(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if walker, ok := map[string]func(*envh.EnvTree, []string) (bool, error){
		"CHYLE_FEATURES":       func(*envh.EnvTree, []string) (bool, error) { return true, nil },
		"CHYLE_GIT_REFERENCE":  c.validateChyleGitReference,
		"CHYLE_GIT_REPOSITORY": c.validateChyleGitRepository,
	}[strings.Join(keyChain, "_")]; ok {
		return walker(fullconfig, keyChain)
	}

	if processor, ok := map[string]func() configurater{
		"CHYLE_DECORATORS_ENV":  func() configurater { return &envDecoratorConfigurator{config: fullconfig} },
		"CHYLE_DECORATORS_JIRA": func() configurater { return jiraDecoratorConfigurator(fullconfig) },
		"CHYLE_EXTRACTORS":      func() configurater { return &extractorsConfigurator{config: fullconfig} },
		"CHYLE_MATCHERS":        func() configurater { return &matchersConfigurator{config: fullconfig} },
		"CHYLE_SENDERS_GITHUB":  func() configurater { return &githubSenderConfigurator{config: fullconfig} },
		"CHYLE_SENDERS_STDOUT":  func() configurater { return &stdoutSenderConfigurator{config: fullconfig} },
	}[strings.Join(keyChain, "_")]; ok {
		return processor().process(c)
	}

	return false, nil
}

func (c *CHYLE) validateChyleGitRepository(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REPOSITORY"}, []string{"PATH"})
}

func (c *CHYLE) validateChyleGitReference(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REFERENCE"}, []string{"FROM", "TO"})
}

func resolveConfig(envConfig *envh.EnvTree) error {
	return envConfig.PopulateStruct(&chyleConfig)
}

func debugConfig() {
	if !EnableDebugging {
		return
	}

	if d, err := json.MarshalIndent(chyleConfig, "", "    "); err == nil {
		logger.Println(string(d))
	}
}
