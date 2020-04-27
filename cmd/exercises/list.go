package exercises

import (
	"encoding/csv"
	"fmt"
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var (
	shouldPrintCsv bool
	includeEmpty bool

	header = []string{"ID", "Kennung", "Nachname", "Vorname", "Abgabe"}
)

var submissionsListCommand = &cobra.Command{
	Use:   "list [exerciseId] [assignmentId]",
	Short: "Lists all submissions within an submissions",
	SilenceErrors: true,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		client := util.NewIliasClient()


		spin := util.StartSpinner("Fetching submissions")
		submissions, err := client.Exercise.List(&ilias.ListParams{
			Reference:    args[0],
			Assignment:   args[1],
			IncludeEmpty: includeEmpty,
		})

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccess(fmt.Sprintf("Received %d entries", len(submissions)))

		if shouldPrintCsv {
			printCsv(submissions)
		} else {
			printTable(submissions)
		}
	},
}

func printCsv(submissions []ilias.SubmissionMeta)  {
	writer := csv.NewWriter(os.Stdout)
	writer.Write(header)

	for _, submission := range submissions {
		writer.Write(submission.ToRow())
	}

	writer.Flush()
}

func printTable(submissions []ilias.SubmissionMeta) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, submission := range submissions {
		table.Append(submission.ToRow())
	}

	table.Render()
}

func init() {
	submissionsListCommand.Flags().BoolVar(&shouldPrintCsv, "csv", false, "Prints the table in csv format")
	submissionsListCommand.Flags().BoolVar(&includeEmpty, "empty", false, "Includes empty submissions")
}
