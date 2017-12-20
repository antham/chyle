package builder

import (
	"bytes"
	"testing"

	"github.com/antham/strumt"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvPrompt(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "1\n"
	p := NewEnvPrompt(EnvConfig{"TEST", "NEXT_TEST", "TEST_NEW_ENV_PROMPT", "Enter a value", func(value string) error { return nil }}, store)

	s := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), &stdout)
	s.AddLinePrompter(p.(strumt.LinePrompter))
	s.SetFirst("TEST")
	s.Run()

	scenario := s.Scenario()

	assert.Len(t, scenario, 1)
	assert.Equal(t, scenario[0].PromptString(), "Enter a value")
	assert.Len(t, scenario[0].Inputs(), 1)
	assert.Equal(t, scenario[0].Inputs()[0], "1")
	assert.Nil(t, scenario[0].Error())

	assert.Equal(t, &Store{"TEST_NEW_ENV_PROMPT": "1"}, store)
}

func TestNewEnvPromptWithEmptyValueGiven(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "\n1\n"
	p := NewEnvPrompt(EnvConfig{"TEST", "NEXT_TEST", "TEST_NEW_ENV_PROMPT", "Enter a value", func(value string) error { return nil }}, store)

	s := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), &stdout)
	s.AddLinePrompter(p.(strumt.LinePrompter))
	s.SetFirst("TEST")
	s.Run()

	scenario := s.Scenario()

	assert.Len(t, scenario, 2)
	assert.Equal(t, scenario[0].PromptString(), "Enter a value")
	assert.Len(t, scenario[0].Inputs(), 1)
	assert.Equal(t, scenario[0].Inputs()[0], "")
	assert.EqualError(t, scenario[0].Error(), "No value given")
	assert.Equal(t, scenario[1].PromptString(), "Enter a value")
	assert.Len(t, scenario[1].Inputs(), 1)
	assert.Equal(t, scenario[1].Inputs()[0], "1")
	assert.Nil(t, scenario[1].Error())

	assert.Equal(t, &Store{"TEST_NEW_ENV_PROMPT": "1"}, store)
}

func TestNewEnvPrompts(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "1\n2\n"
	p := NewEnvPrompts([]EnvConfig{
		{"TEST1", "TEST2", "TEST_PROMPT_1", "Enter a value for prompt 1", func(value string) error { return nil }},
		{"TEST2", "", "TEST_PROMPT_2", "Enter a value for prompt 2", func(value string) error { return nil }},
	}, store)

	s := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), &stdout)
	for _, item := range p {
		switch prompt := item.(type) {
		case strumt.LinePrompter:
			s.AddLinePrompter(prompt)
		case strumt.MultilinePrompter:
			s.AddMultilinePrompter(prompt)
		}
	}
	s.SetFirst("TEST1")
	s.Run()

	scenario := s.Scenario()

	assert.Len(t, scenario, 2)
	assert.Equal(t, scenario[0].PromptString(), "Enter a value for prompt 1")
	assert.Len(t, scenario[0].Inputs(), 1)
	assert.Equal(t, scenario[0].Inputs()[0], "1")
	assert.Nil(t, scenario[0].Error())
	assert.Equal(t, scenario[1].PromptString(), "Enter a value for prompt 2")
	assert.Len(t, scenario[1].Inputs(), 1)
	assert.Equal(t, scenario[1].Inputs()[0], "2")
	assert.Nil(t, scenario[1].Error())

	assert.Equal(t, &Store{"TEST_PROMPT_1": "1", "TEST_PROMPT_2": "2"}, store)
}
