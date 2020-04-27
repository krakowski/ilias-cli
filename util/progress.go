package util

import (
	"github.com/cheggaaa/pb/v3"
	"time"
)

const (
	progressRefreshRate = time.Millisecond * 100
	progressBarWidth = 80
	progressBarTemplate = `{{ (cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" | white) }} {{string . "bar_title"}} {{counters . "%s/%s"}} {{ bar . "[" "=" ">" " " "]"}} {{etime .}} {{percent .}}`

)

func StartProgressBar(total int, title string) (*pb.ProgressBar) {
	progress := pb.New(total)
	progress.SetRefreshRate(progressRefreshRate)
	progress.SetTemplateString(progressBarTemplate)
	progress.SetMaxWidth(progressBarWidth)
	progress.Set("bar_title", title)
	progress.Start()
	return progress
}
