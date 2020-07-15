package grades

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "grades",
	Short: "Provides functions for managing grades",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(exportGradesCommand)
}
