package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var pluginManager *PluginManager

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	pluginManager = &PluginManager{}
	cmds := pluginManager.init()

	rootCmd = &cobra.Command{
		Use:   "utilitybelt",
		Short: "UtilityTool belt core command",
		Long: `UtilityTool belt is an automation CLI for Go that carries out a number of common tasks.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Utility Belt 0.1")
		},
	}

	rootCmd.AddCommand(cmds...)
}
