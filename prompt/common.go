package prompt

import (
	"github.com/antham/strumt"

	"github.com/antham/chyle/prompt/internal/builder"
)

func mergePrompters(prompters ...[]strumt.Prompter) []strumt.Prompter {
	results := prompters[0]

	for _, p := range prompters[1:] {
		results = append(results, p...)
	}

	return results
}

func addMainMenuChoice(choices []builder.SwitchConfig) []builder.SwitchConfig {
	return append(choices, builder.SwitchConfig{"m", "Menu", "mainMenu"})
}

func addQuitChoice(choices []builder.SwitchConfig) []builder.SwitchConfig {
	return append(choices, builder.SwitchConfig{"q", "Dump generated configuration and quit", ""})
}

func addMainMenuAndQuitChoice(choices []builder.SwitchConfig) []builder.SwitchConfig {
	return addQuitChoice(addMainMenuChoice(choices))
}
