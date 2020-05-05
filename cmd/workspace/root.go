package workspace

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "workspace",
	Short: "Provides functions for managing workspaces",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(workspaceInitCommand)
	RootCommand.AddCommand(workspaceUploadCommand)
	RootCommand.AddCommand(workspaceStatusCommand)
}
