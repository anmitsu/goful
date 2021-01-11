// Package widget provides base user interfaces.
package widget

import (
	"sync"

	"github.com/anmitsu/goful/look"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Widget describes a window manager in CUI.
type Widget interface {
	// Viewer and controller
	Draw()            // sets cells to the terminal
	Input(key string) // the key event control

	// Connector
	Next() Widget // retruns a next widget for drawing and input
	Disconnect()  // disconnect reference to the next widget for exiting

	// Window
	Width() int
	Height() int
	LeftTop() (x, y int)
	RightBottom() (x, y int)
	LeftBottom() (x, y int)
	RightTop() (x, y int)
	Resize(x, y, width, height int)
	ResizeRelative(x, y, width, height int)
}

// Window represents a quadrangular on coordinates of x and y.
type Window struct {
	x      int // a horizontal coordinate of the window left top
	y      int // a vertical coordinate of the window left top
	width  int
	height int
}

// NewWindow creates a new window specified coordinates and sizes.
func NewWindow(x, y, width, height int) *Window { return &Window{x, y, width, height} }

// Width returns the window width.
func (w *Window) Width() int { return w.width }

// Height returns the window height.
func (w *Window) Height() int { return w.height }

// LeftTop returns left top coordinates of the window.
func (w *Window) LeftTop() (x, y int) { return w.x, w.y }

// RightBottom returns right bottom coordinates of the window.
func (w *Window) RightBottom() (x, y int) { return w.x + w.width - 1, w.y + w.height - 1 }

// LeftBottom returns left bottom coordinates of the window.
func (w *Window) LeftBottom() (x, y int) { return w.x, w.y + w.height - 1 }

// RightTop returns right top coordinates of the window.
func (w *Window) RightTop() (x, y int) { return w.x + w.width - 1, w.y }

// Border draws a frame with runes represent lines and corners.
func (w *Window) Border() {
	xend, yend := w.RightBottom()
	fg, bg := look.Default().Fg(), look.Default().Bg()
	termbox.SetCell(w.x, w.y, ulCorner, fg, bg)
	for x := w.x + 1; x < xend; x++ {
		termbox.SetCell(x, w.y, hLine, fg, bg)
	}
	termbox.SetCell(xend, w.y, urCorner, fg, bg)

	for y := w.y + 1; y < yend; y++ {
		termbox.SetCell(w.x, y, vLine, fg, bg)
		termbox.SetCell(xend, y, vLine, fg, bg)
	}

	termbox.SetCell(w.x, yend, llCorner, fg, bg)
	for x := w.x + 1; x < xend; x++ {
		termbox.SetCell(x, yend, hLine, fg, bg)
	}
	termbox.SetCell(xend, yend, lrCorner, fg, bg)
}

// Draw the window to cells.
func (w *Window) Draw() {
	w.Clear()
	w.Border()
}

// Clear the window with blanks.
func (w *Window) Clear() {
	xend, yend := w.RightBottom()
	for y := w.y; y < yend+1; y++ {
		for x := w.x; x < xend+1; x++ {
			termbox.SetCell(x, y, ' ', look.Default().Fg(), look.Default().Bg())
		}
	}
}

// Resize the window to coordinates and sizes.
func (w *Window) Resize(x, y, width, height int) {
	w.x, w.y = x, y
	w.width, w.height = width, height
}

// ResizeRelative resizes relative to current sizes.
func (w *Window) ResizeRelative(x, y, width, height int) {
	w.x += x
	w.y += y
	w.width += width
	w.height += height
}

var (
	vLine    rune = 0x2502 // │
	hLine    rune = 0x2500 // ─
	ulCorner rune = 0x250c // ┌
	urCorner rune = 0x2510 // ┐
	llCorner rune = 0x2514 // └
	lrCorner rune = 0x2518 // ┘
)

// SetBorder sets window border runes.
func SetBorder(v, h, ul, ur, ll, lr rune) {
	vLine, hLine, ulCorner, urCorner, llCorner, lrCorner = v, h, ul, ur, ll, lr
}

var mutex sync.Mutex

// Flush is to a terminal refresh with a single thread.
func Flush() {
	mutex.Lock()
	if termbox.IsInit {
		termbox.Flush()
	}
	mutex.Unlock()
}

type (
	// Keymap represents callback functions for the widget input event.
	Keymap map[string]func()
	// Extmap represents callback functions for the widget input event on specified conditions.
	Extmap map[string]map[string]func()
)

var ctrlcombo = map[termbox.Key]string{
	termbox.KeyCtrlSpace:      "C-scape",
	termbox.KeyCtrlA:          "C-a",
	termbox.KeyCtrlB:          "C-b",
	termbox.KeyCtrlC:          "C-c",
	termbox.KeyCtrlD:          "C-d",
	termbox.KeyCtrlE:          "C-e",
	termbox.KeyCtrlF:          "C-f",
	termbox.KeyCtrlG:          "C-g",
	termbox.KeyCtrlH:          "C-h",
	termbox.KeyCtrlI:          "C-i",
	termbox.KeyCtrlJ:          "C-j",
	termbox.KeyCtrlK:          "C-k",
	termbox.KeyCtrlL:          "C-l",
	termbox.KeyCtrlM:          "C-m",
	termbox.KeyCtrlN:          "C-n",
	termbox.KeyCtrlO:          "C-o",
	termbox.KeyCtrlP:          "C-p",
	termbox.KeyCtrlQ:          "C-q",
	termbox.KeyCtrlR:          "C-r",
	termbox.KeyCtrlS:          "C-s",
	termbox.KeyCtrlT:          "C-t",
	termbox.KeyCtrlU:          "C-u",
	termbox.KeyCtrlV:          "C-v",
	termbox.KeyCtrlW:          "C-w",
	termbox.KeyCtrlX:          "C-x",
	termbox.KeyCtrlY:          "C-y",
	termbox.KeyCtrlZ:          "C-z",
	termbox.KeyCtrlLsqBracket: "C-[",
	termbox.KeyCtrlBackslash:  "C-\\",
	termbox.KeyCtrlRsqBracket: "C-]",
	termbox.KeyCtrl6:          "C-6",
	termbox.KeyCtrlSlash:      "C-/",
	termbox.KeySpace:          "space",
	termbox.KeyBackspace2:     "backspace",
}

var specialkey = map[termbox.Key]string{
	termbox.KeyF1:         "f1",
	termbox.KeyF2:         "f2",
	termbox.KeyF3:         "f3",
	termbox.KeyF4:         "f4",
	termbox.KeyF5:         "f5",
	termbox.KeyF6:         "f6",
	termbox.KeyF7:         "f7",
	termbox.KeyF8:         "f8",
	termbox.KeyF9:         "f9",
	termbox.KeyF10:        "f10",
	termbox.KeyF11:        "f11",
	termbox.KeyF12:        "f12",
	termbox.KeyInsert:     "insert",
	termbox.KeyDelete:     "delete",
	termbox.KeyHome:       "home",
	termbox.KeyEnd:        "end",
	termbox.KeyPgup:       "pgup",
	termbox.KeyPgdn:       "pgdn",
	termbox.KeyArrowUp:    "up",
	termbox.KeyArrowDown:  "down",
	termbox.KeyArrowLeft:  "left",
	termbox.KeyArrowRight: "right",
}

// EventToString converts the keyboard intput event to string.
// A meta key input event returns prefixed `M-'.
func EventToString(ev *termbox.Event) string {
	const meta = "M-"
	switch ev.Type {
	case termbox.EventKey:
		if ev.Mod&termbox.ModAlt != 0 {
			if ev.Ch == 0 && ev.Key < 128 {
				return meta + ctrlcombo[ev.Key]
			}
			return meta + string(ev.Ch)
		}
		if ev.Key >= termbox.KeyArrowRight {
			return specialkey[ev.Key]
		} else if ev.Ch == 0 && ev.Key < 128 {
			return ctrlcombo[ev.Key]
		}
		return string(ev.Ch)
	case termbox.EventResize:
		return "resize"
	}
	return ""
}

// SetCells sets a string to cells in a window and returns the last x position.
func SetCells(x, y int, s string, l look.Look) (pos int) {
	pos = x
	for _, r := range s {
		termbox.SetCell(pos, y, r, l.Fg(), l.Bg())
		if runewidth.EastAsianWidth && runewidth.IsAmbiguousWidth(r) {
			termbox.SetCell(pos+1, y, ' ', l.Fg(), l.Bg())
		}
		pos += runewidth.RuneWidth(r)
	}
	return
}

// TruncLeft truncates a string with w cells for the left.
func TruncLeft(s string, w int, head string) string {
	if runewidth.StringWidth(s) <= w {
		return s
	}
	r := []rune(s)
	tw := runewidth.StringWidth(head)
	w -= tw
	width := 0
	i := len(r) - 1
	for ; i >= 0; i-- {
		cw := runewidth.RuneWidth(r[i])
		if width+cw > w {
			break
		}
		width += cw
	}
	return head + string(r[i+1:])
}
