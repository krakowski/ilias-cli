package settings

import (
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

var settingsSynchronizeCommand = &cobra.Command{
	Use:   "synchronize [courseId] [settingsFile]",
	Short: "Lists all members within a course",
	SilenceErrors: true,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			log.Fatal("reading settings failed: " + err.Error())
		}

		settings := ilias.CourseSettings{}
		err = yaml.Unmarshal(data, &settings)
		if err != nil {
			log.Fatal("unmarshalling settings failed: " + err.Error())
		}

		client := util.NewIliasClient()

		spin := util.StartSpinner("Synchronizing settings")
		err = client.Course.SynchronizeSettings(&ilias.SettingsParams{
			Reference: args[0],
			Token: client.User.Token,
		}, &settings)

		if err != nil {
			spin.StopError(err)
			os.Exit(1)
		}

		spin.StopSuccess("Synchronized settings")
	},
}

func init() {

}
