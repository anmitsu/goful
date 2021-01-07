package progbar

import (
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

var progbar *widget.ProgressGauge

// Draw the progress bar.
func Draw() {
	if !progbar.IsFinished() {
		progbar.Draw()
	}
}

// Start the progress bar to an arrival value.
func Start(maxval float64) {
	progbar.Start(maxval)
}

// Finish the progress bar and the clear display.
func Finish() {
	progbar.Finish()
	progbar.Clear()
}

// IsFinished reports whether the progress bar finished.
func IsFinished() bool {
	return progbar.IsFinished()
}

// Update the progress bar by a value.
func Update(value float64) {
	progbar.Update(value)
}

// Resize the progress bar.
func Resize(width, height int) {
	w := int(float64(width) * 0.7)
	x := width - w
	progbar.Resize(x, height-1, w, 1)
	progbar.Clear()
}

// Init initializes the progress bar at the bottom position.
func Init() {
	width, height := termbox.Size()
	w := int(float64(width) * 0.7)
	x := width - w
	progbar = widget.NewProgressGauge(x, height-1, w, 1)
}
