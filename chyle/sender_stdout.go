package chyle

import (
	"encoding/json"
	"fmt"
	"io"
)

// jSONStdoutSender output commit payload as JSON on stdout
type jSONStdoutSender struct {
	stdout io.Writer
}

// Send produces an output on stdout
func (j jSONStdoutSender) Send(changelog *Changelog) error {
	return json.NewEncoder(j.stdout).Encode(changelog)
}

// templateStdoutSender output commit payload using given template on stdout
type templateStdoutSender struct {
	template string
	stdout   io.Writer
}

// Send produces an output on stdout
func (t templateStdoutSender) Send(changelog *Changelog) error {
	datas, err := populateTemplate("stdout-template", t.template, changelog)

	if err != nil {
		return err
	}

	fmt.Fprint(t.stdout, datas)

	return nil
}
