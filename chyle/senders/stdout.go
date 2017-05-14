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

// jSONStdoutSender output commit payload as JSON on stdout
type jSONStdoutSender struct {
	stdout io.Writer
}

// Send produces an output on stdout
func (j jSONStdoutSender) Send(changelog *types.Changelog) error {
	return json.NewEncoder(j.stdout).Encode(changelog)
}

// templateStdoutSender output commit payload using given template on stdout
type templateStdoutSender struct {
	stdout   io.Writer
	template string
}

// Send produces an output on stdout
func (t templateStdoutSender) Send(changelog *types.Changelog) error {
	datas, err := populateTemplate("stdout-template", t.template, changelog)

	if err != nil {
		return err
	}

	fmt.Fprint(t.stdout, datas)

	return nil
}

func buildStdoutSender(config stdoutConfig) Sender {
	if config.FORMAT == "json" {
		return buildJSONStdoutSender()
	}

	return buildTemplateStdoutSender(config.TEMPLATE)
}

func buildJSONStdoutSender() Sender {
	return jSONStdoutSender{
		os.Stdout,
	}
}

func buildTemplateStdoutSender(template string) Sender {
	return templateStdoutSender{
		os.Stdout,
		template,
	}
}
