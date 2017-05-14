package senders

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/antham/chyle/chyle/types"
)

type stdoutConfig struct {
	FORMAT   string
	TEMPLATE string
}

// jSONStdout output commit payload as JSON on stdout
type jSONStdout struct {
	stdout io.Writer
}

// Send produces an output on stdout
func (j jSONStdout) Send(changelog *types.Changelog) error {
	return json.NewEncoder(j.stdout).Encode(changelog)
}

// templateStdout output commit payload using given template on stdout
type templateStdout struct {
	stdout   io.Writer
	template string
}

// Send produces an output on stdout
func (t templateStdout) Send(changelog *types.Changelog) error {
	datas, err := populateTemplate("stdout-template", t.template, changelog)

	if err != nil {
		return err
	}

	fmt.Fprint(t.stdout, datas)

	return nil
}

func buildStdout(config stdoutConfig) Sender {
	if config.FORMAT == "json" {
		return buildJSONStdout()
	}

	return buildTemplateStdout(config.TEMPLATE)
}

func buildJSONStdout() Sender {
	return jSONStdout{
		os.Stdout,
	}
}

func buildTemplateStdout(template string) Sender {
	return templateStdout{
		os.Stdout,
		template,
	}
}
