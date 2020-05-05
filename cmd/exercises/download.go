package exercises

import (
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
)

const (
	SubmissionFilename = "Abgabe"
)

var exerciseDownloadCommand = &cobra.Command{
	Use:   "download [exercise] [assignment] [member]",
	Short: "Downloads a specific submission",
	SilenceErrors: true,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		client := util.NewIliasClient()

		// Add rar mime type
		if err := mime.AddExtensionType(".rar", "application/x-rar"); err != nil {
			log.Fatal(err)
		}

		spin := util.StartSpinner("Downloading submission")

		// Download next submission
		submission, err := client.Exercise.Download(&ilias.DownloadParams{
			Reference:  args[0],
			Assignment: args[1],
			Member:     args[2],
		})

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		// Get file extension
		extensions, err := mime.ExtensionsByType(submission.ContentType)
		if err != nil {
			log.Fatal(err)
		}

		if len(extensions) == 0 {
			log.Fatal("unknown content type ", submission.ContentType)
		}

		downloadPath := filepath.Join(SubmissionFilename + extensions[0])
		err = ioutil.WriteFile(downloadPath, submission.Content, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		spin.StopSuccess(util.NoMessage)
	},
}


