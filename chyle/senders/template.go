package senders

import (
	"bytes"
	"fmt"
	tmpl "html/template"
)

// ErrTemplateMalformed issued when something wrong
// happened when creation a template
type ErrTemplateMalformed struct {
	err error
}

// Error return string error
func (e ErrTemplateMalformed) Error() string {
	return fmt.Sprintf("check your template is well-formed : %s", e.err.Error())
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
