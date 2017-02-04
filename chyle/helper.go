package chyle

import (
	"bytes"
	"fmt"
	tmpl "html/template"
)

func populateTemplate(ID string, template string, data interface{}) (string, error) {
	t := tmpl.New(ID)
	t, err := t.Parse(template)

	if err != nil {
		return "", fmt.Errorf("Check your template is well-formed : %s", err.Error())
	}

	b := bytes.Buffer{}
	err = t.Execute(&b, data)

	if err != nil {
		return "", err
	}

	return b.String(), nil
}
