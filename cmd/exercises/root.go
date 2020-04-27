package exercises

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "exercises",
	Short: "Provides functions for managing exercises",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(submissionsListCommand)
	RootCommand.AddCommand(submissionsDistributeCommand)
}
