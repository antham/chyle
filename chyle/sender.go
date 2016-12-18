package chyle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
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
func CreateSenders(config *viper.Viper) (*[]Sender, error) {
	results := []Sender{}

	for sectionKey := range config.GetStringMap("senders") {
		var ex Sender
		var err error
		switch sectionKey {
		case "stdout":
			if !config.IsSet("senders.stdout.format") {
				err = fmt.Errorf(`"format" key must be defined`)
			}

			ex, err = NewStdoutSender(config.GetString("senders.stdout.format"))
		default:
			err = fmt.Errorf(`"%s" is not a valid sender structure`, sectionKey)
		}

		if err != nil {
			return &[]Sender{}, err
		}

		results = append(results, ex)
	}

	return &results, nil
}
