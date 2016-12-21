package chyle

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPopulateTemplate(t *testing.T) {
	r, err := populateTemplate("test", "{{.test}}", map[string]string{"test": "Hello world !"})

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "Hello world !", r, "Must interpolate data in template")
}

func TestPopulateTemplateWithABadTemplate(t *testing.T) {
	_, err := populateTemplate("test", "{{.test", map[string]string{"test": "Hello world !"})

	assert.EqualError(t, err, "Check your template is well-formed : template: test:1: unclosed action", "Must return an error when template is not well-formed")
}

func TestCheckArguments(t *testing.T) {
	v := viper.New()
	v.Set("test", "test")

	err := checkArguments([]string{"test"}, v)

	assert.NoError(t, err, "Must return no error when data is in config")
}

func TestCheckArgumentsWhenEntryDoesntNotExist(t *testing.T) {
	v := viper.New()

	err := checkArguments([]string{"test"}, v)

	assert.EqualError(t, err, `"test" must be defined`, "Must return an error when config doesn't not exist")
}
