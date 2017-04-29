package chyle

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
func createSenders() *[]sender {
	results := []sender{}

	if chyleConfig.FEATURES.HASGITHUBRELEASESENDER {
		results = append(results, buildGithubReleaseSender())
	}

	if chyleConfig.FEATURES.HASSTDOUTSENDER {
		results = append(results, buildStdoutSender())
	}

	return &results
}
