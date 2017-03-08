package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// sender define where the date must be sent
type sender interface {
	Send(changelog *Changelog) error
}

// Send forward changelog to senders
func Send(senders *[]sender, changelog *Changelog) error {
	for _, sender := range *senders {
		err := sender.Send(changelog)

		if err != nil {
			return err
		}
	}

	return nil
}

// createSenders build senders from a config
func createSenders(config *envh.EnvTree) (*[]sender, error) {
	results := []sender{}

	var se sender
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
			return &[]sender{}, err
		}

		results = append(results, se)
	}

	return &results, nil
}
