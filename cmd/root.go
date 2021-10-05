package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/krakowski/ilias-cli/cmd/exercises"
	"github.com/krakowski/ilias-cli/cmd/grades"
	"github.com/krakowski/ilias-cli/cmd/members"
	"github.com/krakowski/ilias-cli/cmd/workspace"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
)

const (
	version = "1.4.4"
	releaseUrl = "https://api.github.com/repos/krakowski/ilias-cli/releases/latest"
	downloadUrlTemplate = "https://github.com/krakowski/ilias-cli/releases/download/%s/ilias-%s-%s"
)

var rootCommand = &cobra.Command{
	Use:   "ilias",
	Short: "A simple command line interface for managing ILIAS",
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		spin := util.StartSpinner("Checking for new version")

		resp, err := http.Get(releaseUrl)
		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		var release struct {
			Tag  string  `json:"tag_name"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccessRemove()
		if release.Tag != version {
			fmt.Fprintf(os.Stderr, "A new version (%s) is availbale. Please download it using the following url.\n\n", release.Tag)

			var downloadUrl string
			if runtime.GOOS == "linux" {
				downloadUrl = fmt.Sprintf(downloadUrlTemplate, release.Tag, release.Tag, "linux-amd64")
			} else if runtime.GOOS == "darwin" {
				downloadUrl = fmt.Sprintf(downloadUrlTemplate, release.Tag, release.Tag, "darwin-amd64")
			} else if runtime.GOOS == "windows" {
				downloadUrl = fmt.Sprintf(downloadUrlTemplate, release.Tag, release.Tag, "windows-amd64.exe")
			} else {
				log.Fatal("unknown operating system")
			}

			fmt.Fprintf(os.Stderr, "  %s\n\n", downloadUrl)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(exercises.RootCommand)
	rootCommand.AddCommand(members.RootCommand)
	rootCommand.AddCommand(workspace.RootCommand)
	rootCommand.AddCommand(grades.RootCommand)
	rootCommand.Version = version
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
