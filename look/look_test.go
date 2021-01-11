package look

import (
	"testing"

	"github.com/nsf/termbox-go"
)

func TestShowColorAttribute(t *testing.T) {
	if err := termbox.Init(); err != nil {
		t.Fatal(err)
	}
	defer termbox.Close()

	colors := []termbox.Attribute{
		termbox.ColorDefault,
		termbox.ColorBlack,
		termbox.ColorRed,
		termbox.ColorGreen,
		termbox.ColorYellow,
		termbox.ColorBlue,
		termbox.ColorMagenta,
		termbox.ColorCyan,
		termbox.ColorWhite,
		termbox.ColorDarkGray,
		termbox.ColorLightRed,
		termbox.ColorLightGreen,
		termbox.ColorLightYellow,
		termbox.ColorLightBlue,
		termbox.ColorLightMagenta,
		termbox.ColorLightCyan,
		termbox.ColorLightGray,
	}
	attrs := []termbox.Attribute{ // foreground effect attributes
		0,
		termbox.AttrBold,
		termbox.AttrBlink,
		termbox.AttrHidden,
		termbox.AttrDim,
		termbox.AttrUnderline,
		termbox.AttrCursive,
		termbox.AttrReverse, // same effect for background
	}

	// foreground color and attribute combines
	x, y := 0, 0
	offsetX := 0
	for _, attr := range attrs {
		for _, bg := range colors {
			for _, fg := range colors {
				termbox.SetCell(x, y, 'F', fg|attr, bg)
				x++
			}
			x = offsetX
			y++
		}
		offsetX += len(colors) + 1
		x = offsetX
		y = 0
	}

	// background color and attribute combines
	x, y = 0, len(colors)+1
	offsetX = 0
	for _, attr := range attrs {
		for _, bg := range colors {
			for _, fg := range colors {
				termbox.SetCell(x, y, 'B', fg, bg|attr)
				x++
			}
			x = offsetX
			y++
		}
		offsetX += len(colors) + 1
		x = offsetX
		y = len(colors) + 1
	}
	termbox.Flush()
	termbox.PollEvent()
}
