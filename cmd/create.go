package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/chyle/chyle"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new changelog",
	Long: `Create a new changelog according to what is defined as config.

Changelog creation follows this process :

1 - fetch commits
2 - filter relevant commits
3 - extract informations from fetched datas
4 - contact third part services to retrieve additional informations from extracted datas
5 - send result to third part services`,
	Run: func(cmd *cobra.Command, args []string) {
		err := chyle.BuildChangelog(envTree)

		if err != nil {
			failure(err)

			exitError()
		}

		exitSuccess()
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
