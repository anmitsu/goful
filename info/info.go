// Package info displays the file information.
package info

import (
	"os"

	"github.com/anmitsu/goful/widget"
)

var info *infoWindow

// Draw the information bar for file information.
func Draw(fi os.FileInfo) {
	info.draw(fi)
}

// Resize the information bar.
func Resize(x, y, width, height int) {
	info.Resize(x, y, width, height)
}

// ResizeRelative resizes relative to current size.
func ResizeRelative(x, y, width, height int) {
	info.Clear()
	info.ResizeRelative(x, y, width, height)
}

// Init the information bar at the bottom position.
func Init() {
	width, height := widget.Size()
	info = &infoWindow{widget.NewWindow(0, height-1, width, 1)}
}

// infoWindow is the information window to display file information.
type infoWindow struct {
	*widget.Window
}
