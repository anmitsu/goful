// Package widget provides base user interfaces.
package widget

import (
	"github.com/anmitsu/goful/look"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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
	border BorderStyle
}

// NewWindow creates a new window specified coordinates and sizes.
func NewWindow(x, y, width, height int) *Window { return &Window{x, y, width, height, ULBorder} }

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

// Border draws a rectangular frame.
func (w *Window) Border() {
	switch w.border {
	case AllBorder:
		w.BorderUL()
		w.BorderLR()
	case ULBorder:
		w.BorderUL()
	}
}

// BorderUL draws upper and lower frames.
func (w *Window) BorderUL() {
	xend, yend := w.RightBottom()
	for x := w.x; x <= xend; x++ {
		screen.SetContent(x, w.y, hLine, nil, look.Default())  // top side
		screen.SetContent(x, yend, hLine, nil, look.Default()) // bottom side
	}
}

// BorderLR draws left and right frames.
func (w *Window) BorderLR() {
	xend, yend := w.RightBottom()
	for y := w.y + 1; y < yend; y++ {
		screen.SetContent(w.x, y, vLine, nil, look.Default())  // left side
		screen.SetContent(xend, y, vLine, nil, look.Default()) // right side
	}
	screen.SetContent(w.x, w.y, ulCorner, nil, look.Default())
	screen.SetContent(w.x, yend, llCorner, nil, look.Default())
	screen.SetContent(xend, w.y, urCorner, nil, look.Default())
	screen.SetContent(xend, yend, lrCorner, nil, look.Default())
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
			screen.SetContent(x, y, ' ', nil, look.Default())
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

// BorderStyle is a window border style.
type BorderStyle int

const (
	// AllBorder is a style draws all lines and corners.
	AllBorder BorderStyle = iota
	// ULBorder is a style draws upper and lower lines.
	ULBorder
	// NoBorder dose not draw borders.
	NoBorder
)

// SetBorderStyle sets the border style.
func (w *Window) SetBorderStyle(style BorderStyle) {
	w.border = style
}

// BorderStyle returns the border style.
func (w *Window) BorderStyle() BorderStyle {
	return w.border
}

var (
	vLine    rune = tcell.RuneVLine
	hLine    rune = tcell.RuneHLine
	ulCorner rune = tcell.RuneULCorner
	urCorner rune = tcell.RuneURCorner
	llCorner rune = tcell.RuneLLCorner
	lrCorner rune = tcell.RuneLRCorner
)

// SetBorder sets window border runes.
func SetBorder(v, h, ul, ur, ll, lr rune) {
	vLine, hLine, ulCorner, urCorner, llCorner, lrCorner = v, h, ul, ur, ll, lr
}

type (
	// Keymap represents callback functions for the widget input event.
	Keymap map[string]func()
	// Extmap represents callback functions for the widget input event on specified conditions.
	Extmap map[string]map[string]func()
)

var keyToSting = map[tcell.Key]string{
	tcell.KeyCtrlA:         "C-a",
	tcell.KeyCtrlB:         "C-b",
	tcell.KeyCtrlC:         "C-c",
	tcell.KeyCtrlD:         "C-d",
	tcell.KeyCtrlE:         "C-e",
	tcell.KeyCtrlF:         "C-f",
	tcell.KeyCtrlG:         "C-g",
	tcell.KeyCtrlH:         "C-h",
	tcell.KeyCtrlI:         "C-i",
	tcell.KeyCtrlJ:         "C-j",
	tcell.KeyCtrlK:         "C-k",
	tcell.KeyCtrlL:         "C-l",
	tcell.KeyCtrlM:         "C-m",
	tcell.KeyCtrlN:         "C-n",
	tcell.KeyCtrlO:         "C-o",
	tcell.KeyCtrlP:         "C-p",
	tcell.KeyCtrlQ:         "C-q",
	tcell.KeyCtrlR:         "C-r",
	tcell.KeyCtrlS:         "C-s",
	tcell.KeyCtrlT:         "C-t",
	tcell.KeyCtrlU:         "C-u",
	tcell.KeyCtrlV:         "C-v",
	tcell.KeyCtrlW:         "C-w",
	tcell.KeyCtrlX:         "C-x",
	tcell.KeyCtrlY:         "C-y",
	tcell.KeyCtrlZ:         "C-z",
	tcell.KeyCtrlLeftSq:    "C-[",
	tcell.KeyCtrlBackslash: "C-\\",
	tcell.KeyBackspace2:    "backspace",
	tcell.KeyF1:            "f1",
	tcell.KeyF2:            "f2",
	tcell.KeyF3:            "f3",
	tcell.KeyF4:            "f4",
	tcell.KeyF5:            "f5",
	tcell.KeyF6:            "f6",
	tcell.KeyF7:            "f7",
	tcell.KeyF8:            "f8",
	tcell.KeyF9:            "f9",
	tcell.KeyF10:           "f10",
	tcell.KeyF11:           "f11",
	tcell.KeyF12:           "f12",
	tcell.KeyInsert:        "insert",
	tcell.KeyDelete:        "delete",
	tcell.KeyHome:          "home",
	tcell.KeyEnd:           "end",
	tcell.KeyPgUp:          "pgup",
	tcell.KeyPgDn:          "pgdn",
	tcell.KeyUp:            "up",
	tcell.KeyDown:          "down",
	tcell.KeyLeft:          "left",
	tcell.KeyRight:         "right",
}

// EventToString converts the keyboard intput event to string.
// A meta key input event returns prefixed `M-'.
func EventToString(ev *tcell.EventKey) string {
	const meta = "M-"
	if ev.Modifiers() == tcell.ModAlt {
		if ev.Key() < 128 {
			return meta + keyToSting[ev.Key()]
		}
		return meta + string(ev.Rune())
	}
	if key, ok := keyToSting[ev.Key()]; ok {
		return key
	}
	return string(ev.Rune())
}

var screen tcell.Screen

// Init initializes the tcell screen.
func Init() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	} else if err := s.Init(); err != nil {
		panic(err)
	}
	screen = s
}

// Fini finishes the tcell screen.
func Fini() {
	screen.Fini()
}

// Size returns the tcell screen width and height.
func Size() (width, height int) {
	return screen.Size()
}

// ShowCursor shows the cursor at x, y.
func ShowCursor(x, y int) {
	screen.ShowCursor(x, y)
}

// HideCursor hides the cursor.
func HideCursor() {
	screen.HideCursor()
}

// PollEvent polls input events.
func PollEvent() tcell.Event {
	return screen.PollEvent()
}

// Show setted cells.
func Show() {
	screen.Show()
}

// SetCells sets a string to cells in a window and returns the last x position.
func SetCells(x, y int, s string, style tcell.Style) (pos int) {
	pos = x
	for _, r := range s {
		screen.SetContent(pos, y, r, nil, style)
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
