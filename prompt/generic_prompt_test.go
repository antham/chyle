package prompt

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/antham/strumt"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvPrompt(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "1\n"
	p := newEnvPrompt(envConfig{"TEST", "NEXT_TEST", "TEST_NEW_ENV_PROMPT", "Enter a value"}, store)

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
	p := newEnvPrompt(envConfig{"TEST", "NEXT_TEST", "TEST_NEW_ENV_PROMPT", "Enter a value"}, store)

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

func TestNewGroupEnvPromptWithCounter(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "test0\ntest1\n1\ntest2\ntest3\nq\n"
	prompts := newGroupEnvPromptWithCounter(
		[]envConfig{
			envConfig{"TEST_0", "TEST_1", "TEST_*_0", "Enter a value"},
			envConfig{"TEST_1", "choice", "TEST_*_1", "Enter a value"},
		}, store)

	var choice = []strumt.Prompter{
		&switchPrompt{
			"choice",
			[]switchChoice{
				switchChoice{
					"1", "Add new test values", "TEST_0",
				},
				switchChoice{
					"q", "Quit", "",
				},
			},
		},
	}

	prompts = append(prompts, choice...)

	s := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), &stdout)

	for _, item := range prompts {
		switch prompt := item.(type) {
		case strumt.LinePrompter:
			s.AddLinePrompter(prompt)
		case strumt.MultilinePrompter:
			s.AddMultilinePrompter(prompt)
		}
	}

	s.SetFirst("TEST_0")
	s.Run()

	scenario := s.Scenario()

	steps := []struct {
		input string
		err   error
	}{
		{
			"test0",
			nil,
		},
		{
			"test1",
			nil,
		},
		{
			"1",
			nil,
		},
		{
			"test2",
			nil,
		},
		{
			"test3",
			nil,
		},
		{
			"q",
			nil,
		},
	}

	for i, step := range steps {
		assert.Nil(t, step.err)
		assert.Len(t, scenario[i].Inputs(), 1)
		assert.Equal(t, scenario[i].Inputs()[0], step.input)
	}

	assert.Equal(t, &Store{"TEST_0_0": "test0", "TEST_0_1": "test1", "TEST_1_0": "test2", "TEST_1_1": "test3"}, store)
}

func TestNewEnvPrompts(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "1\n2\n"
	p := newEnvPrompts([]envConfig{
		envConfig{"TEST1", "TEST2", "TEST_PROMPT_1", "Enter a value for prompt 1"},
		envConfig{"TEST2", "", "TEST_PROMPT_2", "Enter a value for prompt 2"},
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

func TestNewPromptWithCustomHandlers(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "wrong\nright\n"

	p := newPromptWithCustomHandlers(
		envConfig{"test", "", "TEST", "Enter a value"},
		func(val string) string {
			return ""
		},
		func(err error) string {
			return "test"
		},
		func(val string) error {
			if val != "right" {
				return fmt.Errorf(`Value must be right`)
			}

			(*store)["TEST"] = val

			return nil
		},
		store,
	)

	s := strumt.NewPromptsFromReaderAndWriter(bytes.NewBufferString(buf), &stdout)
	s.AddLinePrompter(p.(strumt.LinePrompter))
	s.SetFirst("test")
	s.Run()

	scenario := s.Scenario()

	assert.Len(t, scenario, 2)
	assert.Equal(t, scenario[0].PromptString(), "Enter a value")
	assert.Len(t, scenario[0].Inputs(), 1)
	assert.Equal(t, scenario[0].Inputs()[0], "wrong")
	assert.EqualError(t, scenario[0].Error(), "Value must be right")
	assert.Equal(t, scenario[1].PromptString(), "Enter a value")
	assert.Len(t, scenario[1].Inputs(), 1)
	assert.Equal(t, scenario[1].Inputs()[0], "right")
	assert.Nil(t, scenario[1].Error())

	assert.Equal(t, &Store{"TEST": "right"}, store)
}
