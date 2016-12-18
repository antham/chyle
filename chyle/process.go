package chyle

import (
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

// process represents all configuration operations defined
// needed to create a changelog
type process struct {
	matchers   *[]Matcher
	extractors *[]Extracter
	expanders  *[]Expander
	senders    *[]Sender
}

// buildProcess creates process entity from defined configuration
func buildProcess(viper *viper.Viper) (*process, error) {
	m, err := CreateMatchers(viper)

	if err != nil {
		return nil, err
	}

	ext, err := CreateExtractors(viper)

	if err != nil {
		return nil, err
	}

	exp, err := CreateExpanders(viper)

	if err != nil {
		return nil, err
	}

	s, err := CreateSenders(viper.GetStringMap("senders"))

	if err != nil {
		return nil, err
	}

	return &process{
		m,
		ext,
		exp,
		s,
	}, nil
}

// proceed extracts datas from a set of commits
func proceed(process *process, commits *[]git.Commit) error {
	comExt, err := Extract(process.extractors, TransformCommitsToMap(Filter(process.matchers, commits)))

	if err != nil {

		return err
	}

	comExp, err := Expand(process.expanders, comExt)

	if err != nil {
		return err
	}

	return Send(process.senders, comExp)
}
