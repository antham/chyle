package tmplh

import (
	"bytes"
	"fmt"
	tmpl "html/template"

	"github.com/Masterminds/sprig"
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

var store = map[string]interface{}{}

func isset(key string) bool {
	_, ok := store[key]

	return ok
}

func set(key string, value interface{}) string {
	store[key] = value

	return ""
}

func get(key string) interface{} {
	return store[key]
}

// Build returns, using a template and generic data, an executed template
func Build(ID string, template string, data interface{}) (string, error) {
	funcMap := sprig.FuncMap()
	funcMap["isset"] = isset
	funcMap["set"] = set
	funcMap["get"] = get

	t, err := tmpl.New(ID).Funcs(funcMap).Parse(template)

	if err != nil {
		return "", ErrTemplateMalformed{err}
	}

	b := bytes.Buffer{}

	if err = t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
