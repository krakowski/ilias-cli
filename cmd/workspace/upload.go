package workspace

import (
	"fmt"
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)


var workspaceUploadCommand = &cobra.Command{
	Use:   "upload",
	Short: "Uploads the current workspace to the ILIAS platform",
	SilenceErrors: true,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		client := util.NewIliasClient()
		workSpace := util.GetWorkspace()
		memberIds := workSpace.Corrections[client.User.Username]

		var corrections []ilias.Correction
		for _, member := range memberIds {
			path := filepath.Join(member, CorrectionFilename)
			correction, err := util.ReadCorrection(path)
			if err != nil {
				log.Fatal(err)
			}

			corrections = append(corrections, *correction)
		}

		// Check if all submissions were corrected
		if filtered := util.FilterCorrections(corrections, func(c ilias.Correction) bool { return !c.Corrected}); len(filtered) > 0 {
			var students []string
			for _, correction := range filtered {
				students = append(students, correction.Student)
			}

			fmt.Fprintf(os.Stderr, "Submissions %v are not corrected yet, aborting", students)
			os.Exit(1)
		}

		// Initialize progress bar
		bar := util.StartProgressBar(len(memberIds), "Uploading corrections")

		// Upload comments
		for _, correction := range corrections {
			err := client.Exercise.UpdateComment(&ilias.CommentParams{
				Reference:  workSpace.Exercise.Reference,
				Assignment: workSpace.Exercise.Assignment,
			}, correction)

			if err != nil {
				log.Fatal(err)
			}

			bar.Increment()
		}

		bar.Finish()

		spin := util.StartSpinner("Updating grades")
		err := client.Exercise.UpdateGrades(&ilias.GradesQuery{
			Reference:  workSpace.Exercise.Reference,
			Assignment: workSpace.Exercise.Assignment,
			Token:      client.User.Token,
		}, corrections)

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccess(fmt.Sprintf("Updated %d entries", len(corrections)))

		spin = util.StartSpinner("Uploading table")
		sheet := util.CreateCorrectionSheet(workSpace.Table.Name, corrections)
		err = client.Tables.Import(&ilias.ImportParams{
			Reference: workSpace.Table.Reference,
			Table:     workSpace.Table.Identifier,
			Token: 	   client.User.Token,
		}, sheet)

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccess(fmt.Sprintf("Uploaded %d entries", len(corrections)))
	},
}
