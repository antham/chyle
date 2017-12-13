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
	Short: "Configuration prompt",
	Run: func(cmd *cobra.Command, args []string) {
		prompts := createPrompt(os.Stdin, os.Stdout)
		store, err := prompts.Run()

		if err != nil {
			failure(err)

			exitError()
		}

		printWithNewLine("")
		printWithNewLine("Generated configuration :")
		printWithNewLine("")

		for key, value := range *store {
			printWithNewLine(fmt.Sprintf("export %s=%s", key, value))
		}
	},
}

var createPrompt = func(reader io.Reader, writer io.Writer) prompt.Prompts {
	return prompt.New(reader, writer)
}

func init() {
	RootCmd.AddCommand(configCmd)
}
