package prompt

import (
	"fmt"
)

type switchChoice struct {
	choice       string
	promptString string
	nextPromptID string
}

type switchPrompt struct {
	iD      string
	choices []switchChoice
}

func (s *switchPrompt) ID() string {
	return s.iD
}

func (s *switchPrompt) PromptString() string {
	out := fmt.Sprintf("Choose one of this option and press enter:\n")

	for _, choice := range s.choices {
		out += fmt.Sprintf("%s - %s\n", choice.choice, choice.promptString)
	}

	return out
}

func (s *switchPrompt) Parse(value string) error {
	if value == "" {
		return fmt.Errorf("No value given")
	}

	for _, choice := range s.choices {
		if choice.choice == value {
			return nil
		}
	}

	return fmt.Errorf("This choice doesn't exist")
}

func (s *switchPrompt) NextOnSuccess(value string) string {
	for _, choice := range s.choices {
		if choice.choice == value {
			return choice.nextPromptID
		}
	}

	return ""
}

func (s *switchPrompt) NextOnError(err error) string {
	return s.iD
}

func (s *switchPrompt) PrintPrompt(prompt string) {
	fmt.Printf("%s", prompt)
}
