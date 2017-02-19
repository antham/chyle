package chyle

import (
	"bytes"
	"fmt"
	tmpl "html/template"
	"log"
	"os"
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

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "CHYLE - ", log.Ldate|log.Ltime)
}

func populateTemplate(ID string, template string, data interface{}) (string, error) {
	t := tmpl.New(ID)
	t, err := t.Parse(template)

	if err != nil {
		return "", ErrTemplateMalformed{err}
	}

	b := bytes.Buffer{}
	err = t.Execute(&b, data)

	if err != nil {
		return "", err
	}

	return b.String(), nil
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

func debug(format string, v ...interface{}) {
	if EnableDebugging {
		logger.Printf(format, v...)
	}
}

// concatErrors transforms an array of error in one error
// by merging error message
func concatErrors(errs *[]error) error {
	if len(*errs) == 0 {
		return nil
	}

	errStr := ""

	for i, e := range *errs {
		errStr += e.Error()

		if i != len(*errs)-1 {
			errStr += ", "
		}
	}

	return fmt.Errorf(errStr)
}
