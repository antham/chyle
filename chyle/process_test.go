package chyle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProcessWithAnEmptyConfig(t *testing.T) {
	chyleConfig = CHYLE{}

	p := buildProcess()

	expected := process{
		&[]matcher{},
		&[]extracter{},
		&map[string][]decorater{},
		&[]sender{},
	}

	assert.EqualValues(t, expected, *p)

}
