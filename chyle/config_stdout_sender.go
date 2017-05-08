package chyle

import (
	"fmt"
	"strings"

	"github.com/antham/envh"
)

// stdoutSenderConfigurator validates stdout sender config defined through environment variables
type stdoutSenderConfigurator struct {
	config *envh.EnvTree
}

func (s *stdoutSenderConfigurator) process(config *CHYLE) (bool, error) {
	if s.isDisabled() {
		return false, nil
	}

	config.FEATURES.HASSENDERS = true
	config.FEATURES.HASSTDOUTSENDER = true

	return false, s.validateFormat()
}

// isDisabled checks if stdout sender is enabled
func (s *stdoutSenderConfigurator) isDisabled() bool {
	return featureDisabled(s.config, [][]string{{"CHYLE", "SENDERS", "STDOUT"}})
}

// validateFormat checks format is a supported stdout format
func (s *stdoutSenderConfigurator) validateFormat() error {
	var err error
	var format string
	keyChain := []string{"CHYLE", "SENDERS", "STDOUT"}

	if format, err = s.config.FindString(append(keyChain, "FORMAT")...); err != nil {
		return ErrMissingEnvVar{[]string{strings.Join(append(keyChain, "FORMAT"), "_")}}
	}

	switch format {
	case "json":
		return nil
	case "template":
		return s.validateTemplateFormat()
	}

	return fmt.Errorf(`"CHYLE_SENDERS_STDOUT_FORMAT" "%s" doesn't exist`, format)
}

// validateTemplateFormat checks a template key is defined
// and template is a valid template
func (s *stdoutSenderConfigurator) validateTemplateFormat() error {
	tmplKeyChain := []string{"CHYLE", "SENDERS", "STDOUT", "TEMPLATE"}

	if ok, err := s.config.HasSubTreeValue(tmplKeyChain...); !ok || err != nil {
		return ErrMissingEnvVar{[]string{strings.Join(tmplKeyChain, "_")}}
	}

	if err := validateTemplate(s.config, tmplKeyChain); err != nil {
		return err
	}

	return nil
}
