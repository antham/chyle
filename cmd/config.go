package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/antham/chyle/prompt"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate environments variables from a prompt session",
	Run: func(cmd *cobra.Command, args []string) {
		prompts := createPrompt(os.Stdin, os.Stdout)
		store, err := prompts.Run()

		if err != nil {
			failure(err)

			exitError()
		}

		printWithNewLine("")
		printWithNewLine("Generated environments variables :")
		printWithNewLine("")

		for key, value := range *store {
			printWithNewLine(fmt.Sprintf("%s=%s", key, value))
		}
	},
}

var createPrompt = func(reader io.Reader, writer io.Writer) prompt.Prompts {
	return prompt.New(reader, writer)
}

func init() {
	RootCmd.AddCommand(configCmd)
}
