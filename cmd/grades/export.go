package grades

import (
    "encoding/csv"
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/krakowski/ilias"
    "github.com/krakowski/ilias-cli/util"
)

var exportGradesCommand = &cobra.Command{
    Use:   "export [exerciseId]",
    Short: "Exports all grades within an exercise",
    SilenceErrors: true,
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {

        client := util.NewIliasClient()


        spin := util.StartSpinner("Fetching submissions")
        grades, err := client.Exercise.Export(&ilias.GradesExportQuery{
            Reference:args[0],
        })

        if err != nil {
            spin.StopError(err)
            os.Exit(1)
        }

        spin.StopSuccess(fmt.Sprintf("Received %d entries", len(grades)))

        printCsv(grades)
    },
}

func printCsv(grades []ilias.Grading)  {
    writer := csv.NewWriter(os.Stdout)
    writer.Write(grades[0].ToHeader())

    for _, grading := range grades {
        writer.Write(grading.ToRow())
    }

    writer.Flush()
}
