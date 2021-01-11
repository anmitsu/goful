// Package infobar displays the file infomation and file control task.
package infobar

import (
	"fmt"
	"os"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var infobar *InfoBar

// Draw the infomation bar for file infomation.
func Draw(fi os.FileInfo) {
	infobar.Clear()
	if infobar.task != nil {
		infobar.drawTask()
	} else {
		infobar.draw(fi)
	}
}

// StartTask starts drawing the file control task.
func StartTask(fi os.FileInfo) {
	infobar.task = fi
	infobar.drawTask()
}

// FinishTask finishes drawing the file control task.
func FinishTask() {
	infobar.done++
}

// ExistsTask reports whether the exists task.
func ExistsTask() bool { return infobar.task != nil }

// StartTaskCount starts the task count.
func StartTaskCount(count int) {
	infobar.done = 0
	infobar.taskCount = count
}

// ResetTaskCount resets the task count.
func ResetTaskCount() {
	infobar.done = 0
	infobar.taskCount = 0
	infobar.task = nil
	infobar.Clear()
}

// Resize the infomation bar.
func Resize(x, y, width, height int) {
	infobar.Resize(x, y, width, height)
}

// ResizeRelative resizes relative to current size.
func ResizeRelative(x, y, width, height int) {
	infobar.Clear()
	infobar.ResizeRelative(x, y, width, height)
}

// Init the infomation bar at the bottom position.
func Init() {
	width, height := termbox.Size()
	infobar = &InfoBar{widget.NewWindow(0, height-1, width, 1), nil, 0, 0}
}

// InfoBar is the infomation bar to display file infomation and the file control task.
type InfoBar struct {
	*widget.Window
	task      os.FileInfo
	taskCount int
	done      int
}

func (w *InfoBar) drawTask() {
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
