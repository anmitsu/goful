// Package infobar is the infomation bar to display file infomation.
package infobar

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var infobar *InfoBar

// Draw the infomation bar for file infomation.
func Draw(fi os.FileInfo) {
	infobar.Clear()
	if infobar.ctrling != nil {
		infobar.drawFilectrl()
	} else {
		infobar.draw(fi)
	}
}

// StartFilectrl starts drawing file controls for fi.
func StartFilectrl(fi os.FileInfo) {
	infobar.ctrling = fi
	Draw(fi)
}

// FinishFilectrl finishes file controls drawing.
func FinishFilectrl() {
	infobar.ctrling = nil
	infobar.Clear()
}

// Resize the infomation bar.
func Resize(width, height int) {
	infobar.Resize(0, height-1, width, 1)
	infobar.Clear()
}

// Init the infomation bar at the bottom left half position.
func Init() {
	width, height := termbox.Size()
	infobar = &InfoBar{widget.NewWindow(0, height-1, width, 1), nil}
}

// InfoBar is the infomation bar to display file infomation.
type InfoBar struct {
	*widget.Window
	ctrling os.FileInfo
}

func (w *InfoBar) drawFilectrl() {
	const (
		Gb = 1024 * 1024 * 1024
		Mb = 1024 * 1024
		kb = 1024
	)
	size := w.ctrling.Size()
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
	width := int(float64(w.Width()) * 0.3)
	name := filepath.Base(w.ctrling.Name())
	s := fmt.Sprintf("Processing (%s): %s", format, name)
	s = runewidth.Truncate(s, width, "...")
	widget.SetCells(x, y, s, look.Default())
}
