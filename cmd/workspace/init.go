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

		workspace := util.ReadWorkspace()
		client := util.NewIliasClient()

		// Parse workspace file
		memberIds := workspace.Corrections[client.User.Username]

		// Add rar mime type
		if err := mime.AddExtensionType(".rar", "application/x-rar"); err != nil {
			log.Fatal(err)
		}

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
				Reference:  workspace.Exercise.Reference,
				Assignment: workspace.Exercise.Assignment,
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

			if len(extensions) == 0 {
				log.Fatal("unknown content type ", submission.ContentType)
			}

			// Save submission
			downloadPath := filepath.Join(memberId, SubmissionFilename + extensions[0])
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


