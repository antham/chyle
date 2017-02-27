package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// Sender define where the date must be sent
type Sender interface {
	Send(*[]map[string]interface{}) error
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
func CreateSenders(config *envh.EnvTree) (*[]Sender, error) {
	results := []Sender{}

	var se Sender
	var subConfig envh.EnvTree
	var err error

	for _, k := range config.GetChildrenKeys() {
		switch k {
		case "STDOUT":
			subConfig, err = config.FindSubTree("STDOUT")

			if err != nil {
				break
			}

			se, err = buildStdoutSender(&subConfig)
		case "GITHUB":
			subConfig, err = config.FindSubTree("GITHUB")

			if err != nil {
				break
			}

			se, err = buildGithubReleaseSender(&subConfig)
		default:
			err = fmt.Errorf(`a wrong sender key containing "%s" was defined`, k)
		}

		if err != nil {
			return &[]Sender{}, err
		}

		results = append(results, se)
	}

	return &results, nil
}
