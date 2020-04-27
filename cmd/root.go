package cmd

import (
	"github.com/krakowski/ilias-cli/cmd/members"
	"github.com/krakowski/ilias-cli/cmd/exercises"
	"github.com/krakowski/ilias-cli/cmd/workspace"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "ilias",
	Short: "A simple command line interface for managing ILIAS",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(exercises.RootCommand)
	rootCommand.AddCommand(members.RootCommand)
	rootCommand.AddCommand(workspace.RootCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
