package workspace

import (
	"fmt"
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"log"
	"os"
)


var workspaceUploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "Uploads the current workspace to the ILIAS platform",
	SilenceErrors: true,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		workspace := util.ReadWorkspace()
		client := util.NewIliasClient()

		members := workspace.Corrections[client.User.Username]
		corrections := util.ReadCorrections(members)

		// Check if all submissions were corrected
		stats := util.GetCorrectionStats(corrections)
		if len(stats.Pending) > 0 {
			fmt.Fprintln(os.Stderr, util.Red("not all submissions are corrected yet"))
			os.Exit(1)
		}

		// Initialize progress bar
		bar := util.StartProgressBar(len(members), "Uploading corrections")

		// Upload comments
		for _, correction := range corrections {
			err := client.Exercise.UpdateComment(&ilias.CommentParams{
				Reference:  workspace.Exercise.Reference,
				Assignment: workspace.Exercise.Assignment,
			}, correction)

			if err != nil {
				log.Fatal(err)
			}

			bar.Increment()
		}

		bar.Finish()

		spin := util.StartSpinner("Updating grades")
		err := client.Exercise.UpdateGrades(&ilias.GradesQuery{
			Reference:  workspace.Exercise.Reference,
			Assignment: workspace.Exercise.Assignment,
			Token:      client.User.Token,
		}, corrections)

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccess(fmt.Sprintf("Updated %d entries", len(corrections)))

		//spin = util.StartSpinner("Uploading table")
		//sheet := util.CreateCorrectionSheet(workspace.Table.Name, corrections)
		//err = client.Tables.Import(&ilias.ImportParams{
		//	Reference: workspace.Table.Reference,
		//	Table:     workspace.Table.Identifier,
		//	Token:     client.User.Token,
		//}, sheet)
		//
		//if err != nil {
		//	spin.StopError(err)
		//	os.Exit(1)
		//}
		//
		//spin.StopSuccess(fmt.Sprintf("Uploaded %d entries", len(corrections)))
	},
}
