package prompt

import (
	"github.com/antham/strumt"
)

func mergePrompters(prompters ...[]strumt.Prompter) []strumt.Prompter {
	results := prompters[0]

	for _, p := range prompters[1:] {
		results = append(results, p...)
	}

	return results
}

func addMainMenuChoice(choices []switchChoice) []switchChoice {
	return append(choices, switchChoice{"m", "Menu", "mainMenu"})
}

func addQuitChoice(choices []switchChoice) []switchChoice {
	return append(choices, switchChoice{"q", "Dump generated configuration and quit", ""})
}

func addMainMenuAndQuitChoice(choices []switchChoice) []switchChoice {
	return addQuitChoice(addMainMenuChoice(choices))
}
