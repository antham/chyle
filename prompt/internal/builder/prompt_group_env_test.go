package builder

import (
	"bytes"
	"testing"

	"github.com/antham/strumt"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupEnvPromptWithCounter(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "test0\ntest1\n1\ntest2\ntest3\nq\n"
	prompts := NewGroupEnvPromptWithCounter(
		[]EnvConfig{
			{"TEST_0", "TEST_1", "TEST_*_0", "Enter a value", func(value string) error { return nil }},
			{"TEST_1", "choice", "TEST_*_1", "Enter a value", func(value string) error { return nil }},
		}, store)

	var choice = []strumt.Prompter{
		&switchPrompt{
			"choice",
			[]SwitchConfig{
				{
					"1", "Add new test values", "TEST_0",
				},
				{
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
