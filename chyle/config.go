package chyle

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antham/envh"
)

// validater must be implemented to add a validator when settings
// struct fields
type validater interface {
	validate() (bool, error)
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
		"CHYLE_MATCHERS":             c.validateChyleMatchers,
		"CHYLE_SENDERS_STDOUT":       c.validateChyleSendersStdout,
	}[strings.Join(keyChain, "_")]; ok {
		return walker(fullconfig, keyChain)
	}

	if validator, ok := map[string]func() validater{
		"CHYLE_DECORATORS_JIRA": func() validater { return jiraDecoratorValidator{fullconfig} },
		"CHYLE_SENDERS_GITHUB":  func() validater { return githubSenderValidator{fullconfig} },
	}[strings.Join(keyChain, "_")]; ok {
		return validator().validate()
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

func (c *CHYLE) validateChyleMatchers(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if featureDisabled(fullconfig, [][]string{{"CHYLE", "MATCHERS"}}) {
		return true, nil
	}

	c.MATCHERS = map[string]string{}

	for _, key := range []string{"MESSAGE", "COMMITTER", "AUTHOR"} {
		value, err := fullconfig.FindString("CHYLE", "MATCHERS", key)

		if err != nil {
			continue
		}

		if err := validateRegexp(fullconfig, []string{"CHYLE", "MATCHERS", key}); err != nil {
			return true, err
		}

		c.MATCHERS[key] = value
	}

	value, err := fullconfig.FindString("CHYLE", "MATCHERS", "TYPE")

	if err != nil {
		return true, nil
	}

	if err := validateOneOf(fullconfig, []string{"CHYLE", "MATCHERS", "TYPE"}, []string{regularTypeMatcher, mergeTypeMatcher}); err != nil {
		return true, err
	}

	c.MATCHERS["TYPE"] = value

	return true, nil
}

func (c *CHYLE) validateChyleSendersStdout(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if featureDisabled(fullconfig, [][]string{{"CHYLE", "SENDERS", "STDOUT"}}) {
		return false, nil
	}

	var err error
	var format string

	if format, err = fullconfig.FindString(append(keyChain, "FORMAT")...); err != nil {
		return false, ErrMissingEnvVar{[]string{strings.Join(append(keyChain, "FORMAT"), "_")}}
	}

	switch format {
	case "json":
		return false, nil
	case "template":
		tmplKeyChain := append(keyChain, "TEMPLATE")

		if ok, err := fullconfig.HasSubTreeValue(tmplKeyChain...); !ok || err != nil {
			return false, ErrMissingEnvVar{[]string{strings.Join(tmplKeyChain, "_")}}
		}

		if err := validateTemplate(fullconfig, tmplKeyChain); err != nil {
			return false, err
		}
	default:
		return false, fmt.Errorf(`"CHYLE_SENDERS_STDOUT_FORMAT" "%s" doesn't exist`, format)
	}

	return false, nil
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
