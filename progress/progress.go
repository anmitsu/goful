// Package progress displays the file control task progress.
package progress

import (
	"fmt"
	"os"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var progress *progressWindow

// Draw the progress task and gauge.
func Draw() {
	if !progress.gauge.IsFinished() {
		progress.drawTask()
		progress.gauge.Draw()
	}
}

// Start progressing to an arrival value.
func Start(maxval float64) {
	progress.gauge.Start(maxval)
}

// Finish progressing and the clear display.
func Finish() {
	progress.gauge.Finish()
	progress.Clear()
}

// IsFinished reports whether progressing finished.
func IsFinished() bool {
	return progress.gauge.IsFinished()
}

// Update progressing by a value.
func Update(value float64) {
	progress.gauge.Update(value)
}

// Resize the progress window and gauge.
func Resize(x, y, width, height int) {
	progress.Resize(x, y, width, height)
	progress.gauge.Resize(x, y+1, width, height)
}

// Init initializes the progress window and gauge at the bottom position.
func Init() {
	width, height := termbox.Size()
	progress = &progressWindow{
		widget.NewWindow(0, height-4, width, 1),
		widget.NewProgressGauge(0, height-3, width, 1),
		nil,
		0,
		0,
	}
}

// StartTask starts the file control task.
func StartTask(fi os.FileInfo) {
	progress.task = fi
}

// FinishTask finishes the file control task.
func FinishTask() {
	progress.done++
}

// StartTaskCount starts the task count.
func StartTaskCount(count int) {
	progress.done = 0
	progress.taskCount = count
}

type progressWindow struct {
	*widget.Window
	gauge     *widget.ProgressGauge
	task      os.FileInfo
	taskCount int
	done      int
}

func (w *progressWindow) drawTask() {
	if w.task == nil {
		return
	}
	w.Clear()
	const (
		Gb = 1024 * 1024 * 1024
		Mb = 1024 * 1024
		kb = 1024
	)
	size := w.task.Size()
	format := ""
	if size > Gb {
		format = fmt.Sprintf("%.1fG", float64(size)/Gb)
	} else if size > Mb {
		format = fmt.Sprintf("%.1fM", float64(size)/Mb)
	} else if size > kb {
		format = fmt.Sprintf("%.1fk", float64(size)/kb)
	} else {
		format = fmt.Sprintf("%d", size)
	}

	x, y := w.LeftTop()
	x++
	name := w.task.Name()
	s := fmt.Sprintf("Progress %d/%d (%s): %s", w.done+1, w.taskCount, format, name)
	s = runewidth.Truncate(s, w.Width(), "...")
	widget.SetCells(x, y, s, look.Default())
}
