package chyle

import (
	"bytes"
	tmpl "html/template"
	"log"
	"os"
)

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

func debug(format string, v ...interface{}) {
	if EnableDebugging {
		logger.Printf(format, v...)
	}
}
