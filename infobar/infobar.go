// Package infobar displays the file information.
package infobar

import (
	"os"

	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

var infobar *infoWindow

// Draw the information bar for file information.
func Draw(fi os.FileInfo) {
	infobar.draw(fi)
}

// Resize the information bar.
func Resize(x, y, width, height int) {
	infobar.Resize(x, y, width, height)
}

// ResizeRelative resizes relative to current size.
func ResizeRelative(x, y, width, height int) {
	infobar.Clear()
	infobar.ResizeRelative(x, y, width, height)
}

// Init the information bar at the bottom position.
func Init() {
	width, height := termbox.Size()
	infobar = &infoWindow{widget.NewWindow(0, height-1, width, 1)}
}

// infoWindow is the information window to display file information.
type infoWindow struct {
	*widget.Window
}
