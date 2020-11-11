package workspace

import (
	"errors"
	"fmt"
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"github.com/gabriel-vasile/mimetype"
	"io/ioutil"
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

		// Initialize progress bar
		spin := util.StartSpinner(fmt.Sprintf("Downloading submissions (0/%d)", len(memberIds)))
		for index, memberId := range memberIds {

			correctionPath := filepath.Join(memberId, CorrectionFilename)

			// Skip submissions already downloaded
			if _, err := os.Stat(memberId); !os.IsNotExist(err) {
				spin.UpdateMessage(fmt.Sprintf("Downloading submissions (%d/%d)", index + 1, len(memberIds)))
				continue
			}

			// Download next submission
			submission, err := client.Exercise.Download(&ilias.DownloadParams{
				Reference:  workspace.Exercise.Reference,
				Assignment: workspace.Exercise.Assignment,
				Member:     memberId,
			})

			if err != nil {
				spin.StopError(err)
				os.Exit(1)
			}

			// Detect content type
			mime := mimetype.Detect(submission.Content)
			if len(mime.Extension()) == 0 {
				spin.StopError(errors.New("detecting content type failed"))
				os.Exit(1)
			}

			// Ensure submission directory is present
			if _, err := os.Stat(memberId); os.IsNotExist(err) {
				_ = os.Mkdir(memberId, os.ModePerm)
			}

			// Save submission
			downloadPath := filepath.Join(memberId, SubmissionFilename + mime.Extension())
			err = ioutil.WriteFile(downloadPath, submission.Content, os.ModePerm)
			if err != nil {
				spin.StopError(err)
				os.Exit(1)
			}

			err = util.WriteCorrectionTemplate(correctionPath, util.TemplateParams{
				Student: memberId,
				Tutor:   client.User.Username,
			})

			if err != nil {
				spin.StopError(err)
				os.Exit(1)
			}

			spin.UpdateMessage(fmt.Sprintf("Downloading submissions (%d/%d)", index + 1, len(memberIds)))
		}

		spin.StopSuccess(util.NoMessage)
	},
}


