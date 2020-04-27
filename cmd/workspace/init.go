package workspace

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
	CorrectionFilename = "Korrektur.yml"
)

var workspaceInitCommand = &cobra.Command{
	Use:   "init",
	Short: "Initializes a workspace for corrections",
	SilenceErrors: true,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		client := util.NewIliasClient()

		// Parse workspace file
		workSpace := util.GetWorkspace()
		memberIds := workSpace.Corrections[client.User.Username]

		// Initialize progress bar
		bar := util.StartProgressBar(len(memberIds), "Downloading submissions")
		for _, memberId := range memberIds {

			correctionPath := filepath.Join(memberId, CorrectionFilename)

			// Skip submissions already downloaded
			if _, err := os.Stat(memberId); !os.IsNotExist(err) {
				bar.Increment()
				continue
			}

			// Download next submission
			submission, err := client.Exercise.Download(&ilias.DownloadParams{
				Reference:  workSpace.Exercise.Reference,
				Assignment: workSpace.Exercise.Assignment,
				Member:     memberId,
			})

			if err != nil {
				log.Fatal(err)
			}

			// Get file extension
			extensions, err := mime.ExtensionsByType(submission.ContentType)
			if err != nil {
				log.Fatal(err)
			}

			// Ensure submission directory is present
			if _, err := os.Stat(memberId); os.IsNotExist(err) {
				os.Mkdir(memberId, os.ModePerm)
			}

			// Save submission
			downloadPath := filepath.Join(memberId, SubmissionFilename+ extensions[0])
			err = ioutil.WriteFile(downloadPath, submission.Content, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			err = util.WriteCorrectionTemplate(correctionPath, util.TemplateParams{
				Student: memberId,
				Tutor:   client.User.Username,
			})

			if err != nil {
				log.Fatal(err)
			}

			bar.Increment()
		}

		bar.Finish()
	},
}


