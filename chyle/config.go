package chyle

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/antham/envh"
)

var chyleConfig CHYLE

// ErrMissingEnvVar is triggered when a required environment variable is missing
type ErrMissingEnvVar struct {
	keys []string
}

// Error output error as string
func (e ErrMissingEnvVar) Error() string {
	switch len(e.keys) {
	case 1:
		return fmt.Sprintf(`environment variable missing : "%s"`, e.keys[0])
	default:
		return fmt.Sprintf(`environments variables missing : "%s"`, strings.Join(e.keys, `", "`))
	}
}

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
	}
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
		"CHYLE_DECORATORS_JIRA":      c.validateChyleJiraDecorators,
		"CHYLE_DECORATORS_JIRA_KEYS": c.setJiraKeys,
		"CHYLE_EXTRACTORS":           c.validateChyleExtractors,
		"CHYLE_FEATURES":             c.setFeatures,
		"CHYLE_GIT_REFERENCE":        c.validateChyleGitReference,
		"CHYLE_GIT_REPOSITORY":       c.validateChyleGitRepository,
		"CHYLE_MATCHERS":             c.validateChyleMatchers,
		"CHYLE_SENDERS_GITHUB":       c.validateChyleSendersGithub,
		"CHYLE_SENDERS_STDOUT":       c.validateChyleSendersStdout,
	}[strings.Join(keyChain, "_")]; ok {
		return walker(fullconfig, keyChain)
	}

	return false, nil
}

func (c *CHYLE) checkSubConfigPool(conf *envh.EnvTree, prefix []string, keys []string) error {
	undefinedKeys := []string{}

	for _, key := range keys {
		ok, err := conf.HasSubTreeValue(append(prefix, key)...)

		if !ok || err != nil {
			undefinedKeys = append(undefinedKeys, strings.Join(append(prefix, key), "_"))
		}
	}

	if len(undefinedKeys) > 0 {
		return ErrMissingEnvVar{undefinedKeys}
	}

	return nil
}

func (c *CHYLE) validateURL(fullconfig *envh.EnvTree, chain []string) error {
	if _, err := url.ParseRequestURI(fullconfig.FindStringUnsecured(chain...)); err != nil {
		return fmt.Errorf(`provide a valid URL for "%s", "%s" given`, strings.Join(chain, "_"), fullconfig.FindStringUnsecured(chain...))
	}

	return nil
}

func (c *CHYLE) validateRegexp(fullconfig *envh.EnvTree, chain []string) error {
	if _, err := regexp.Compile(fullconfig.FindStringUnsecured(chain...)); err != nil {
		return fmt.Errorf(`provide a valid regexp for "%s", "%s" given`, strings.Join(chain, "_"), fullconfig.FindStringUnsecured(chain...))
	}

	return nil
}

func (c *CHYLE) validateOneOf(fullconfig *envh.EnvTree, chain []string, choices []string) error {
	val := fullconfig.FindStringUnsecured(chain...)

	for _, choice := range choices {
		if choice == val {
			return nil
		}
	}

	return fmt.Errorf(`provide a value for "%s" from one of those values : ["%s"], "%s" given`, strings.Join(chain, "_"), strings.Join(choices, `", "`), val)
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
	return false, c.checkSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REPOSITORY"}, []string{"PATH"})
}

func (c *CHYLE) validateChyleGitReference(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	return false, c.checkSubConfigPool(fullconfig, []string{"CHYLE", "GIT", "REFERENCE"}, []string{"FROM", "TO"})
}

func (c *CHYLE) validateAndSetChyleDecoratorsEnv(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	keys, err := fullconfig.FindChildrenKeys("CHYLE", "DECORATORS", "ENV")

	if err == nil && len(keys) == 0 {
		return true, nil
	}

	for _, key := range keys {
		if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "DECORATORS", "ENV", key}, []string{"DESTKEY", "FIELD"}); err != nil {
			return true, err
		}
	}

	c.DECORATORS.ENV = map[string]map[string]string{}

	for _, key := range keys {
		c.DECORATORS.ENV[key] = map[string]string{}

		for _, field := range []string{"DESTKEY", "FIELD"} {
			c.DECORATORS.ENV[key][field] = fullconfig.FindStringUnsecured("CHYLE", "DECORATORS", "ENV", key, field)
		}
	}

	return true, nil
}

func (c *CHYLE) validateChyleExtractors(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	keys, err := fullconfig.FindChildrenKeys("CHYLE", "EXTRACTORS")

	if err == nil && len(keys) == 0 {
		return true, nil
	}

	for _, key := range keys {
		if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "EXTRACTORS", key}, []string{"ORIGKEY", "DESTKEY", "REG"}); err != nil {
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

			if err := c.validateRegexp(fullconfig, chain); err != nil {
				return err
			}

			c.EXTRACTORS[key][field] = value
		}
	}

	return nil
}

func (c *CHYLE) validateChyleMatchers(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if !fullconfig.IsExistingSubTree("CHYLE", "MATCHERS") {
		return true, nil
	}

	c.MATCHERS = map[string]string{}

	for _, key := range []string{"MESSAGE", "COMMITTER", "AUTHOR"} {
		value, err := fullconfig.FindString("CHYLE", "MATCHERS", key)

		if err != nil {
			continue
		}

		if err := c.validateRegexp(fullconfig, []string{"CHYLE", "MATCHERS", key}); err != nil {
			return true, err
		}

		c.MATCHERS[key] = value
	}

	value, err := fullconfig.FindString("CHYLE", "MATCHERS", "TYPE")

	if err != nil {
		return true, nil
	}

	if err := c.validateOneOf(fullconfig, []string{"CHYLE", "MATCHERS", "TYPE"}, []string{regularTypeMatcher, mergeTypeMatcher}); err != nil {
		return true, err
	}

	c.MATCHERS["TYPE"] = value

	return true, nil
}

func (c *CHYLE) validateChyleSendersStdout(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if !fullconfig.IsExistingSubTree("CHYLE", "SENDERS", "STDOUT") {
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
		if err = c.checkSubConfigPool(fullconfig, []string{"CHYLE", "SENDERS", "STDOUT"}, []string{"TEMPLATE"}); err != nil {
			return false, err
		}
	default:
		return false, fmt.Errorf(`"CHYLE_SENDERS_STDOUT_FORMAT" "%s" doesn't exist`, format)
	}

	return false, nil
}

func (c *CHYLE) validateChyleSendersGithub(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	if !fullconfig.IsExistingSubTree("CHYLE", "SENDERS", "GITHUB") {
		return false, nil
	}

	if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "SENDERS", "GITHUB", "CREDENTIALS"}, []string{"OAUTHTOKEN", "OWNER"}); err != nil {
		return false, err
	}

	if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "SENDERS", "GITHUB", "RELEASE"}, []string{"TAGNAME", "TEMPLATE"}); err != nil {
		return false, err
	}

	if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "SENDERS", "GITHUB", "REPOSITORY"}, []string{"NAME"}); err != nil {
		return false, err
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

func (c *CHYLE) validateChyleJiraDecorators(fullconfig *envh.EnvTree, keyChain []string) (bool, error) {
	hasCred := fullconfig.IsExistingSubTree("CHYLE", "DECORATORS", "JIRA")
	hasExt := fullconfig.IsExistingSubTree("CHYLE", "EXTRACTORS", "JIRAISSUEID")

	if !hasCred && !hasExt {
		return false, nil
	}

	if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS"}, []string{"URL", "USERNAME", "PASSWORD"}); err != nil {
		return false, err
	}

	if err := c.validateURL(fullconfig, []string{"CHYLE", "DECORATORS", "JIRA", "CREDENTIALS", "URL"}); err != nil {
		return false, err
	}

	keys, err := fullconfig.FindChildrenKeys("CHYLE", "DECORATORS", "JIRA", "KEYS")

	if err != nil {
		return false, fmt.Errorf(`define at least one environment variable couple "CHYLE_DECORATORS_JIRA_KEYS_*_DESTKEY" and "CHYLE_DECORATORS_JIRA_KEYS_*_FIELD", replace "*" with your own naming`)
	}

	for _, key := range keys {
		if err := c.checkSubConfigPool(fullconfig, []string{"CHYLE", "DECORATORS", "JIRA", "KEYS", key}, []string{"DESTKEY", "FIELD"}); err != nil {
			return false, err
		}
	}

	return false, c.checkSubConfigPool(fullconfig, []string{"CHYLE", "EXTRACTORS", "JIRAISSUEID"}, []string{"ORIGKEY", "DESTKEY", "REG"})
}

func resolveConfig(envConfig *envh.EnvTree) error {
	return envConfig.PopulateStruct(&chyleConfig)
}
