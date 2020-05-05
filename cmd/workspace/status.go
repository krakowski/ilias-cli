package workspace

import (
	"fmt"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"os"
)


var workspaceStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Shows information about the current workspace",
	SilenceErrors: true,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		workSpace := util.ReadWorkspace()
		members := workSpace.Corrections[util.ReadUserCache()]
		corrections := util.ReadCorrections(members)

		stats := util.GetCorrectionStats(corrections)
		fmt.Fprintf(os.Stderr, "Inside ILIAS workspace '%s'\n\n", workSpace.Title)

		fmt.Fprintf(os.Stderr, "  %s corrected (%d): ", util.Green(util.IconSuccess), len(stats.Corrected))
		for _, correction := range stats.Corrected {
			fmt.Fprintf(os.Stderr, "%s ", correction.Student)
		}

		fmt.Fprintf(os.Stderr, "\n\n")

		fmt.Fprintf(os.Stderr, "  %s pending (%d): ", util.Red(util.IconError), len(stats.Pending))
		for _, correction := range stats.Pending {
			fmt.Fprintf(os.Stderr, "%s ", correction.Student)
		}

		fmt.Fprintf(os.Stderr, "\n\n")
	},
}
