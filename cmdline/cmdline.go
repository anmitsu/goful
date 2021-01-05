// Package cmdline is a command line widget like shell.
package cmdline

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// A Mode describes cmdline mode.
type Mode interface {
	String() string // name to use as commands
	Prompt() string // displayed head of cmdline
	Result() string
	Init(*Cmdline)
	Draw(*Cmdline)
	Run(*Cmdline)
}

// Cmdline is the text box with specified mode.
type Cmdline struct {
	*widget.TextBox
	filer      widget.Widget
	completion widget.Widget
	mode       Mode
	History    *History
}

var keymap func(*Cmdline) widget.Keymap

// Config sets keymap function for cmdline.
func Config(config func(*Cmdline) widget.Keymap) {
	keymap = config
}

// New creates the cmdline with specified mode and history list box.
// These widget size based on filer of parent widget.
func New(m Mode, filer widget.Widget) *Cmdline {
	x, y := filer.LeftTop()
	width, height := filer.Width(), filer.Height()
	filer.ResizeRelative(0, 0, 0, -1)
	c := &Cmdline{
		TextBox:    widget.NewTextBox(x, y+height-1, width, 1),
		filer:      filer,
		completion: nil,
		mode:       m,
		History:    &History{},
	}
	c.mode.Init(c)

	y = (y + height) * 2 / 3
	height -= y + 1
	c.History = NewHistory(x, y, width, height, c)
	c.Edithook = func() {
		c.History.update()
		c.History.MoveTop()
	}
	c.History.update()
	return c
}

// Next implements widget.Widget.
func (c *Cmdline) Next() widget.Widget { return c.completion }

// Disconnect implements widget.Widget.
func (c *Cmdline) Disconnect() { c.completion = nil }

// StartCompletion starts completion based on cmdline text.
func (c *Cmdline) StartCompletion() {
	x, y := c.History.LeftTop()
	width, height := c.History.Width(), c.History.Height()
	comp := NewCompletion(x, y, width, height, c)
	if comp.IsEmpty() {
		return
	} else if len(comp.List()) == 1 {
		comp.InsertCompletion()
		return
	}
	c.completion = comp
}

// Resize the cmdline and history list box.
func (c *Cmdline) Resize(x, y, width, height int) {
	c.TextBox.Resize(x, y+height-1, width, 1)
	y = (y + height) * 2 / 3
	height -= y + 1
	c.History.Resize(x, y, width, height)
}

// ResizeRelative relative resizes the cmdline and history list box.
func (c *Cmdline) ResizeRelative(x, y, width, height int) {
	c.TextBox.ResizeRelative(x, y, width, height)
	c.History.ResizeRelative(x, y, width, height)
}

// DrawLine draws the cmdline prompt and text.
func (c *Cmdline) DrawLine() {
	c.Clear()
	x, y := c.LeftTop()
	x++
	x = widget.SetCells(x, y, c.mode.Prompt(), look.Prompt())
	w := c.Width() - runewidth.StringWidth(c.mode.Prompt()) - 2
	s := c.String()
	s = runewidth.Truncate(s, w, "")
	if c.Cursor() >= w {
		s = c.TextBeforeCursor()
		s = widget.TruncLeft(s, w, "...")
		x = widget.SetCells(x, y, s, look.Cmdline())
		termbox.SetCursor(x, y)
	} else {
		widget.SetCells(x, y, s, look.Cmdline())
		termbox.SetCursor(x+c.Cursor(), y)
	}
}

// Draw the cmdline mode and completion or histry list box
func (c *Cmdline) Draw() {
	c.mode.Draw(c)
	if c.Next() != nil {
		c.Next().Draw()
	} else {
		c.History.Draw()
	}
}

// Input to text or widget callback function.
func (c *Cmdline) Input(key string) {
	if c.completion != nil {
		c.completion.Input(key)
	} else if cb, ok := keymap(c)[key]; ok {
		cb()
	} else {
		if key == "space" {
			c.InsertChar(' ')
		} else if utf8.RuneCountInString(key) == 1 {
			r, _ := utf8.DecodeRuneInString(key)
			c.InsertChar(r)
		}
	}
}

// Exit the cmdline and add cmdline text to history.
func (c *Cmdline) Exit() {
	c.History.add()
	c.History = nil
	termbox.HideCursor()
	c.filer.ResizeRelative(0, 0, 0, 1)
	c.filer.Disconnect()
}

// Run the cmdline based on mode and add cmdline text to history.
func (c *Cmdline) Run() {
	c.History.add()
	c.mode.Run(c)
}

var historyMap = map[string][]string{}

// LoadHistory loads from the file and append to history map of key as file name.
func LoadHistory(path string) error {
	path = utils.ExpandPath(path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	history := make([]string, 0)
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		history = append(history, string(line))
	}
	key := filepath.Base(path)
	historyMap[key] = history
	return nil
}

// SaveHistory saves to the file of history.
func SaveHistory(path string) error {
	path = utils.ExpandPath(path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	key := filepath.Base(path)
	if history, ok := historyMap[key]; ok {
		writer := bufio.NewWriter(file)
		for _, h := range history {
			if _, err := writer.WriteString(h + "\n"); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// History is the cmdline mode history.
type History struct {
	*widget.ListBox
	cmdline *Cmdline
	input   string // in the input
}

// NewHistory creates the history list box.
func NewHistory(x, y, width, height int, cmdline *Cmdline) *History {
	lb := widget.NewListBox(x, y, width, height, "History")
	lb.SetLower(-1)
	lb.SetCursor(lb.Lower())
	return &History{lb, cmdline, ""}
}

func (h *History) update() {
	text := h.cmdline.String()
	name := h.cmdline.mode.String()

	h.input = text

	if history, ok := historyMap[name]; ok {
		h.ClearList()
		for i := len(history) - 1; i > -1; i-- {
			hist := history[i]
			if strings.Contains(hist, text) {
				h.AppendHighlightString(hist, text)
			}
		}
	}
}

func (h *History) add() {
	text := h.cmdline.String()
	mode := h.cmdline.mode.String()

	if text == "" || text[0] == ' ' {
		return
	}

	if history, ok := historyMap[mode]; ok {
		for i, hist := range history {
			if hist == text {
				history = append(history[:i], history[i+1:]...)
				break
			}
		}
		historyMap[mode] = append(history, text)
	} else {
		history := []string{}
		historyMap[mode] = append(history, text)
	}
}

// Delete history on the cursor.
func (h *History) Delete() {
	if h.Cursor() < 0 || h.Upper() < 1 {
		return
	}
	mode := h.cmdline.mode.String()
	if history, ok := historyMap[mode]; ok {
		for i := 0; i < len(history); i++ {
			if history[i] == h.CurrentContent().Name() {
				history = append(history[:i], history[i+1:]...)
				historyMap[mode] = history

				h.cmdline.SetText(h.input)
				h.update()
				h.AdjustCursor()
				break
			}
		}
	}
}

// MoveCursor moves list box cursor and sets text to the cmdline.
func (h *History) MoveCursor(amount int) {
	h.ListBox.MoveCursor(amount)
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}

// CursorDown down list box cursor and sets text to the cmdline.
func (h *History) CursorDown() {
	h.ListBox.CursorDown()
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}

// CursorUp up list box cursor and sets text to the cmdline.
func (h *History) CursorUp() {
	h.ListBox.CursorUp()
	if h.Cursor() == h.Lower() {
		h.cmdline.SetText(h.input)
	} else {
		h.cmdline.SetText(h.CurrentContent().Name())
	}
}
