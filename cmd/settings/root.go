package settings

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "settings",
	Short: "Provides functions for synchronizing course settings",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(settingsSynchronizeCommand)
}
