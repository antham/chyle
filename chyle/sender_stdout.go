package chyle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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
