package cmd

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	var code int
	var err error
	var w sync.WaitGroup

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	failure = func(e error) {
		err = e
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				code = r.(int)
			}

			w.Done()
		}()

		os.Args = []string{"", "whatever"}

		Execute()
	}()

	w.Wait()

	assert.EqualError(t, err, `unknown command "whatever" for "chyle"`)
	assert.EqualValues(t, 1, code, "Must exit with an errors (exit 1)")
}
