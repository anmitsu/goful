// Package infobar is the infomation bar to display file infomation.
package infobar

import (
	"fmt"
	"os"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/progbar"
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
	Draw(fi)
}

// FinishTask finishes drawing the file control task.
func FinishTask() {
	infobar.task = nil
	infobar.Clear()
}

// Resize the infomation bar.
func Resize(x, y, width, height int) {
	infobar.Resize(x, y, width, height)
}

// Init the infomation bar at the bottom left half position.
func Init() {
	width, height := termbox.Size()
	infobar = &InfoBar{widget.NewWindow(0, height-1, width, 1), nil}
}

// InfoBar is the infomation bar to display file infomation and the file control task.
type InfoBar struct {
	*widget.Window
	task os.FileInfo
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
	s := fmt.Sprintf("Progress (%s): %s | ", format, name)
	width := runewidth.StringWidth(s)
	if width > w.Width()/2 {
		width = w.Width() / 2
	}
	progbar.Resize(width, y, w.Width()-width, w.Height())
	s = runewidth.Truncate(s, width, "...")
	widget.SetCells(x, y, s, look.Default())
}
