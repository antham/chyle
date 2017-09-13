package builder

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/antham/strumt"
	"github.com/stretchr/testify/assert"
)

func TestNewPromptWithCustomHandlers(t *testing.T) {
	store := &Store{}

	var stdout bytes.Buffer
	buf := "wrong\nright\n"

	p := NewPromptWithCustomHandlers(
		EnvConfig{"test", "", "TEST", "Enter a value"},
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
