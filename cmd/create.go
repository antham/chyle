package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/antham/chyle/chyle"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new changelog",
	Long: `Create a new changelog according to what is defined in config file.

Changelog creation follows this process :

1 - fetch commits
2 - filter relevant commits
3 - extract informations from fetched datas
4 - contact third part services to retrieve additional informations from extracted datas
5 - send result to third part services`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := os.Getwd()

		if err != nil {
			fmt.Println(err)

			return
		}

		if len(args) < 2 {
			fmt.Println("Must provides 2 arguments")

			return
		}

		err = chyle.BuildChangelog(path, envTree, args[0], args[1])

		if err != nil {
			fmt.Println(err)

			return
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
