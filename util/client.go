package util

import (
	"fmt"
	"github.com/krakowski/ilias"
	"os"
)

func NewIliasClient() (*ilias.Client) {
	credentials := GetCredentials()
	spin := StartSpinner("Logging in at 'https://ilias.hhu.de'")
	client, err := ilias.NewClient(nil, credentials)
	if err != nil {
		spin.StopError(err)
		os.Exit(1)
	}

	spin.StopSuccess(fmt.Sprintf("Welcome %s!", client.User.Firstname))
	return client
}
