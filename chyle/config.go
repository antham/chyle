package chyle

import (
	"encoding/json"
	"strings"

	"github.com/antham/envh"
)

// configurater must be implemented to process custom config
type configurater interface {
	process() (bool, error)
}

var chyleConfig CHYLE

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

// Walk traverses struct to populate or validate fields
func (c *CHYLE) Walk(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if walker, ok := map[string]func(*envh.EnvTree, []string) (bool, error){
		"CHYLE_DECORATORS_ENV":       c.validateAndSetChyleDecoratorsEnv,
		"CHYLE_DECORATORS_JIRA_KEYS": c.setJiraKeys,
		"CHYLE_EXTRACTORS":           c.validateChyleExtractors,
		"CHYLE_FEATURES":             c.setFeatures,
		"CHYLE_GIT_REFERENCE":        c.validateChyleGitReference,
		"CHYLE_GIT_REPOSITORY":       c.validateChyleGitRepository,
	}[strings.Join(keyChain, "_")]; ok {
		return walker(fullconfig, keyChain)
	}

	if processor, ok := map[string]func() configurater{
		"CHYLE_DECORATORS_JIRA": func() configurater { return jiraDecoratorProcessor{fullconfig} },
		"CHYLE_SENDERS_GITHUB":  func() configurater { return githubSenderProcessor{fullconfig} },
		"CHYLE_SENDERS_STDOUT":  func() configurater { return stdoutSenderProcessor{fullconfig} },
		"CHYLE_MATCHERS":        func() configurater { return &matchersConfigurator{chyleConfig: c, config: fullconfig} },
	}[strings.Join(keyChain, "_")]; ok {
		return processor().process()
	}

	return false, nil
}

func (c *CHYLE) setFeatures(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	structs := []struct {
		ref   *bool
		chain []string
	}{
		{
			&(c.FEATURES.HASDECORATORS),
			[]string{"CHYLE", "DECORATORS"},
		},
		{
			&(c.FEATURES.HASEXTRACTORS),
			[]string{"CHYLE", "EXTRACTORS"},
		},
		{
			&(c.FEATURES.HASMATCHERS),
			[]string{"CHYLE", "MATCHERS"},
		},
		{
			&(c.FEATURES.HASSENDERS),
			[]string{"CHYLE", "SENDERS"},
		},
		{
			&(c.FEATURES.HASJIRADECORATOR),
			[]string{"CHYLE", "DECORATORS", "JIRA"},
		},
		{
			&(c.FEATURES.HASENVDECORATOR),
			[]string{"CHYLE", "DECORATORS", "ENV"},
		},
		{
			&(c.FEATURES.HASGITHUBRELEASESENDER),
			[]string{"CHYLE", "SENDERS", "GITHUB"},
		},
		{
			&(c.FEATURES.HASSTDOUTSENDER),
			[]string{"CHYLE", "SENDERS", "STDOUT"},
		},
	}

	for _, s := range structs {
		if fullconfig.IsExistingSubTree(s.chain...) {
			*(s.ref) = true
		} else {
			*(s.ref) = false
		}
	}

	return true, nil
}

func (c *CHYLE) validateChyleGitRepository(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REPOSITORY"}, []string{"PATH"})
}

func (c *CHYLE) validateChyleGitReference(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, validateSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REFERENCE"}, []string{"FROM", "TO"})
}

func (c *CHYLE) validateAndSetChyleDecoratorsEnv(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if featureDisabled(fullconfig, [][]string{{"CHYLE", "DECORATORS", "ENV"}}) {
		return true, nil
	}

	c.DECORATORS.ENV = map[string]map[string]string{}

	for _, key := range fullconfig.FindChildrenKeysUnsecured("CHYLE", "DECORATORS", "ENV") {
		c.DECORATORS.ENV[key] = map[string]string{}

		if err := validateSubConfigPool(fullconfig, []string{"CHYLE", "DECORATORS", "ENV", key}, []string{"DESTKEY", "VARNAME"}); err != nil {
			return true, err
		}

		for _, field := range []string{"DESTKEY", "VARNAME"} {
			c.DECORATORS.ENV[key][field] = fullconfig.FindStringUnsecured("CHYLE", "DECORATORS", "ENV", key, field)
		}
	}

	return true, nil
}

func (c *CHYLE) validateChyleExtractors(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if featureDisabled(fullconfig, [][]string{{"CHYLE", "EXTRACTORS"}}) {
		return true, nil
	}

	for _, key := range fullconfig.FindChildrenKeysUnsecured("CHYLE", "EXTRACTORS") {
		if err := validateSubConfigPool(fullconfig, []string{"CHYLE", "EXTRACTORS", key}, []string{"ORIGKEY", "DESTKEY", "REG"}); err != nil {
			return true, err
		}
	}

	return true, c.setChyleExtractors(fullconfig, keyChain)
}

func (c *CHYLE) setChyleExtractors(fullconfig *envh.EnvTree, keyChain []string) error {
	c.EXTRACTORS = map[string]map[string]string{}

	for _, key := range fullconfig.FindChildrenKeysUnsecured(keyChain...) {
		c.EXTRACTORS[key] = map[string]string{}

		for _, field := range []string{"ORIGKEY", "DESTKEY", "REG"} {
			chain := []string{"CHYLE", "EXTRACTORS", key, field}

			value := fullconfig.FindStringUnsecured(chain...)

			if err := validateRegexp(fullconfig, chain); err != nil {
				return err
			}

			c.EXTRACTORS[key][field] = value
		}
	}

	return nil
}

func (c *CHYLE) setJiraKeys(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	keys := fullconfig.FindChildrenKeysUnsecured(keyChain...)

	c.DECORATORS.JIRA.KEYS = map[string]string{}

	for _, key := range keys {
		datas := map[string]string{}

		for _, field := range []string{"DESTKEY", "FIELD"} {
			datas[field] = fullconfig.FindStringUnsecured(append(keyChain, key, field)...)
		}

		c.DECORATORS.JIRA.KEYS[datas["DESTKEY"]] = datas["FIELD"]
	}

	return true, nil
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
