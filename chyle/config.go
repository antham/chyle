package chyle

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/antham/envh"

	"github.com/antham/chyle/chyle/decorators"
	"github.com/antham/chyle/chyle/senders"
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
		HASMATCHERS             bool
		HASEXTRACTORS           bool
		HASDECORATORS           bool
		HASSENDERS              bool
		HASJIRAISSUEDECORATOR   bool
		HASGITHUBISSUEDECORATOR bool
		HASENVDECORATOR         bool
		HASGITHUBRELEASESENDER  bool
		HASSTDOUTSENDER         bool
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
	EXTRACTORS map[string]struct {
		ORIGKEY string
		DESTKEY string
		REG     *regexp.Regexp
	}
	DECORATORS decorators.Config
	SENDERS    senders.Config
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
		"CHYLE_DECORATORS_ENV":         func() configurater { return &envDecoratorConfigurator{fullconfig} },
		"CHYLE_DECORATORS_JIRAISSUE":   func() configurater { return jiraIssueDecoratorConfigurator(fullconfig) },
		"CHYLE_DECORATORS_GITHUBISSUE": func() configurater { return githubIssueDecoratorConfigurator(fullconfig) },
		"CHYLE_EXTRACTORS":             func() configurater { return &extractorsConfigurator{fullconfig} },
		"CHYLE_MATCHERS":               func() configurater { return &matchersConfigurator{fullconfig} },
		"CHYLE_SENDERS_GITHUBRELEASE":  func() configurater { return &githubReleaseSenderConfigurator{fullconfig} },
		"CHYLE_SENDERS_STDOUT":         func() configurater { return &stdoutSenderConfigurator{fullconfig} },
	}[strings.Join(keyChain, "_")]; ok {
		return processor().process(c)
	}

	return false, nil
}

func (c *CHYLE) validateChyleGitRepository(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateEnvironmentVariablesDefinition(fullconfig, [][]string{{"CHYLE", "GIT", "REPOSITORY", "PATH"}})
}

func (c *CHYLE) validateChyleGitReference(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateEnvironmentVariablesDefinition(fullconfig, [][]string{{"CHYLE", "GIT", "REFERENCE", "FROM"}, {"CHYLE", "GIT", "REFERENCE", "TO"}})
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
