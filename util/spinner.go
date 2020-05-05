package util

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"time"
)

type Spinner struct {
	spinner 	*spinner.Spinner
}

const (
	IconSuccess = "✓"
	IconError   = "✘"

	NoMessage = ""
)

func StartSpinner(title string) *Spinner {
	spin := spinner.New(spinner.CharSets[14], 100 * time.Millisecond, spinner.WithWriter(os.Stderr))
	spin.Suffix = " " + title
	spin.HideCursor = true
	spin.Start()
	return &Spinner{spinner: spin}
}

func (spin *Spinner) StopSuccess(message string) {
	spin.spinner.Stop()

	if message == NoMessage {
		fmt.Fprintf(os.Stderr, "%s%s\n", Green(IconSuccess), spin.spinner.Suffix)
	} else {
		fmt.Fprintf(os.Stderr, "%s%s - %s\n", Green(IconSuccess), spin.spinner.Suffix, message)
	}
}

func (spin *Spinner) StopError(err error) {
	spin.spinner.Stop()
	fmt.Fprintf(os.Stderr, "%s%s - %s\n", Red(IconError), spin.spinner.Suffix, Red("(" + err.Error() + ")"))
}
