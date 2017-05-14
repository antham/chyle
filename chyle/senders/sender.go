package senders

import (
	"github.com/antham/chyle/chyle/types"
)

// Sender define where the date must be sent
type Sender interface {
	Send(changelog *types.Changelog) error
}

// Send forward changelog to senders
func Send(senders *[]Sender, changelog *types.Changelog) error {
	for _, sender := range *senders {
		err := sender.Send(changelog)

		if err != nil {
			return err
		}
	}

	return nil
}

// CreateSenders build senders from a config
func CreateSenders(features map[string]bool, senders Config) *[]Sender {
	results := []Sender{}

	if enabled, ok := features["githubReleaseSender"]; ok && enabled {
		results = append(results, buildGithubReleaseSender(senders.GITHUBRELEASE))
	}

	if enabled, ok := features["stdoutSender"]; ok && enabled {
		results = append(results, buildStdoutSender(senders.STDOUT))
	}

	return &results
}
