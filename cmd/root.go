package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/chyle/chyle"
	"github.com/antham/envh"
)

var envTree *envh.EnvTree
var debug bool

// RootCmd represents initial cobra command
var RootCmd = &cobra.Command{
	Use:   "chyle",
	Short: "Create a changelog from your commit history",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		failure(err)
		exitError()
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debugging")
}

func initConfig() {
	e, err := envh.NewEnvTree("CHYLE", "_")

	if err != nil {
		failure(err)
		exitError()
	}

	envTree = &e

	chyle.EnableDebugging = debug
}
