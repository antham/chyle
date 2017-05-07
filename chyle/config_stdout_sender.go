package chyle

import (
	"fmt"
	"strings"

	"github.com/antham/envh"
)

// stdoutSenderValidator validates stdout sender config defined through environment variables
type stdoutSenderValidator struct {
	config *envh.EnvTree
}

func (s stdoutSenderValidator) validate() (bool, error) {
	if s.isDisabled() {
		return false, nil
	}

	return false, s.validateFormat()
}

// isDisabled checks if stdout sender is enabled
func (s stdoutSenderValidator) isDisabled() bool {
	return featureDisabled(s.config, [][]string{{"CHYLE", "SENDERS", "STDOUT"}})
}

// validateFormat checks format is a supported stdout format
func (s stdoutSenderValidator) validateFormat() error {
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
func (s stdoutSenderValidator) validateTemplateFormat() error {
	tmplKeyChain := []string{"CHYLE", "SENDERS", "STDOUT", "TEMPLATE"}

	if ok, err := s.config.HasSubTreeValue(tmplKeyChain...); !ok || err != nil {
		return ErrMissingEnvVar{[]string{strings.Join(tmplKeyChain, "_")}}
	}

	if err := validateTemplate(s.config, tmplKeyChain); err != nil {
		return err
	}

	return nil
}
