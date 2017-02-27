package chyle

import (
	"fmt"
	"strings"

	"github.com/antham/envh"
)

type strConfigMapping struct {
	keyChain  []string
	dstValue  *string
	mandatory bool
}

type boolConfigMapping struct {
	keyChain  []string
	dstValue  *bool
	mandatory bool
}

func extractStringConfig(config *envh.EnvTree, mapping []strConfigMapping, prefix []string) error {
	var fullKey, v string
	var err error

	for _, e := range mapping {
		if !e.mandatory && !config.IsExistingSubTree(e.keyChain...) {
			continue
		}

		fullKey = strings.Join(append(prefix, e.keyChain...), "_")
		v, err = config.FindString(e.keyChain...)

		if err != nil {
			return fmt.Errorf(`missing "%s"`, fullKey)
		}

		debug(`"%s" defined with value "%s"`, fullKey, v)

		*(e.dstValue) = v
	}

	return nil
}

func extractBoolConfig(config *envh.EnvTree, mapping []boolConfigMapping, prefix []string) error {
	var fullKey string
	var v bool
	var err error

	for _, e := range mapping {
		if !e.mandatory && !config.IsExistingSubTree(e.keyChain...) {
			continue
		}

		fullKey = strings.Join(append(prefix, e.keyChain...), "_")
		v, err = config.FindBool(e.keyChain...)

		if err != nil {
			return fmt.Errorf(`missing "%s"`, fullKey)
		}

		debug(`"%s" defined with value "%s"`, fullKey, v)

		*(e.dstValue) = v
	}

	return nil
}
