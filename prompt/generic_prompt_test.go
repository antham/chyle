package prompt

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

}
