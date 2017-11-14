package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/antham/chyle/prompt"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var code int
	var wg sync.WaitGroup
	var output []string

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	printWithNewLine = func(str string) {
		output = append(output, str)
	}

	wg.Add(1)

	rd := bytes.NewBufferString("test\ntest\ntest\nq\n")
	wr := bytes.Buffer{}

	createPrompt = func(reader io.Reader, writer io.Writer) prompt.Prompts {
		return prompt.New(rd, &wr)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				code = r.(int)
			}

			wg.Done()
		}()

		os.Args = []string{"", "config"}

		Execute()
	}()

	wg.Wait()

	promptRecord, err := ioutil.ReadAll(&wr)

	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 0, code, "Must exit with no errors (exit 0)")
	assert.Equal(t, "Enter a git commit ID that start your range : \n\nEnter a git commit ID that end your range : \n\nEnter your git path repository : \n\nChoose one of this option and press enter:\n1 - Add a matcher\n2 - Add an extractor\n3 - Add a decorator\n4 - Add a sender\nq - Dump generated configuration and quit\n : \n", string(promptRecord))

	expected := []string{"", "Generated environments variables :", "", "CHYLE_GIT_REFERENCE_FROM=test", "CHYLE_GIT_REFERENCE_TO=test", "CHYLE_GIT_REPOSITORY_PATH=test"}

	sort.Strings(output)
	sort.Strings(expected)

	assert.Equal(t, expected, output)
}
