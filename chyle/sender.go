package chyle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Sender define where the date must be sent
type Sender interface {
	Send(*[]map[string]interface{}) error
}

// StdoutSender send output on stdout
type StdoutSender struct {
	format string
	stdout io.Writer
}

// NewStdoutSender creates a StdoutSender
func NewStdoutSender(format string) (StdoutSender, error) {
	switch format {
	case "json":
		return StdoutSender{
			format,
			os.Stdout,
		}, nil
	default:
		return StdoutSender{}, fmt.Errorf("\"%s\" format does not exist", format)
	}
}

// Send commitMaps to stdout using format
func (s StdoutSender) Send(commitMaps *[]map[string]interface{}) error {
	switch s.format {
	case "json":
		return json.NewEncoder(s.stdout).Encode(commitMaps)
	}

	return nil
}

// Send forward informations extracted before to senders
func Send(senders *[]Sender, commitMaps *[]map[string]interface{}) error {
	for _, sender := range *senders {
		err := sender.Send(commitMaps)

		if err != nil {
			return err
		}
	}

	return nil
}

// CreateSenders build senders from a config
func CreateSenders(senders map[string]interface{}) (*[]Sender, error) {
	results := []Sender{}

	for dk, dv := range senders {
		var ex Sender
		var err error

		e, ok := dv.(map[string]interface{})

		if !ok {
			return &[]Sender{}, fmt.Errorf(`sender "%s" must contains key=value string values`, dk)
		}

		switch dk {
		case "stdout":
			if v, ok := e["format"]; ok {
				s, ok := v.(string)

				if !ok {
					return &[]Sender{}, fmt.Errorf(`extractor "%s" is not a string`, s)
				}

				ex, err = NewStdoutSender(s)
			} else {
				err = fmt.Errorf(`"format" key must be defined`)
			}
		default:
			err = fmt.Errorf(`"%s" is not a valid sender structure`, dk)
		}

		if err != nil {
			return &[]Sender{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
