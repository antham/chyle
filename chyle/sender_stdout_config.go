package chyle

import (
	"fmt"
	"os"

	"github.com/antham/envh"
)

func buildStdoutSender(config *envh.EnvTree) (Sender, error) {
	format, err := config.FindString("FORMAT")

	if err != nil {
		return nil, fmt.Errorf(`missing "SENDERS_STDOUT_FORMAT"`)
	}

	switch format {
	case "json":
		return buildJSONStdoutSender()
	case "template":
		return buildTemplateStdoutSender(config)
	}

	return nil, fmt.Errorf("\"%s\" format does not exist", format)
}

func buildJSONStdoutSender() (Sender, error) {
	return jSONStdoutSender{
		os.Stdout,
	}, nil
}

func buildTemplateStdoutSender(config *envh.EnvTree) (Sender, error) {
	template, err := config.FindString("TEMPLATE")

	if err != nil {
		return nil, fmt.Errorf(`"SENDERS_STDOUT_TEMPLATE" must be defined when "template" format is defined`)
	}

	return templateStdoutSender{
		template,
		os.Stdout,
	}, nil
}
