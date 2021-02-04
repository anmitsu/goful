// Package message provides the message window widget.
package message

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/util"
	"github.com/anmitsu/goful/widget"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Info notifies an information message.
func Info(s string) {
	messenger.info(s)
}

// Infof notifies a formated information message.
func Infof(format string, a ...interface{}) {
	messenger.info(fmt.Sprintf(format, a...))
}

// Error notifies an error message.
func Error(e error) {
	messenger.error(e)
}

// Errorf notifies a formated error message.
func Errorf(format string, a ...interface{}) {
	messenger.error(fmt.Errorf(format, a...))
}

// Resize the message window.
func Resize(x, y, width, height int) {
	messenger.Resize(x, y, width, height)
}

// ResizeRelative resizes relative to current size.
func ResizeRelative(x, y, width, height int) {
	messenger.ResizeRelative(x, y, width, height)
}

// Draw the message window.
func Draw() {
	if !messenger.display {
		messenger.Clear()
	}
}

// Sec sets a display second for the message.
func Sec(sec time.Duration) {
	messenger.sec = sec
}

// SetInfoLog sets a path to log information messages.
// If sets to "", non logging.
func SetInfoLog(path string) {
	messenger.infolog = util.ExpandPath(path)
}

// SetErrorLog sets a path to log error messages.
// If sets to "", non logging.
func SetErrorLog(path string) {
	messenger.errlog = util.ExpandPath(path)
}

// Init initializes the message works asynchronously.
func Init() {
	width, height := widget.Size()
	messenger = &message{
		Window:  widget.NewWindow(0, height-2, width, 1),
		buf:     make(chan buffer, 20),
		display: false,
		sec:     5,
		infolog: "",
		errlog:  "",
	}
	go messenger.run()
}

var messenger *message

type buffer struct {
	string
	tcell.Style
}

type message struct {
	*widget.Window
	buf     chan buffer
	display bool
	sec     time.Duration
	infolog string
	errlog  string
}

func (m *message) info(s string) {
	if m.infolog != "" {
		if err := m.log(m.infolog, s); err != nil {
			m.add(err.Error(), look.MessageError())
		}
	}
	m.add(s, look.MessageInfo())
}

func (m *message) error(e error) {
	if m.errlog != "" {
		if err := m.log(m.errlog, e.Error()); err != nil {
			m.add(err.Error(), look.MessageError())
		}
	}
	m.add(e.Error(), look.MessageError())
}

func (m *message) log(path, s string) error {
	logtime := time.Now().Format("2006-01-02 15:04:30")
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, fmt.Sprintf("%s %s\n", logtime, s))
	if err != nil {
		return err
	}
	return nil
}

func (m *message) add(msg string, s tcell.Style) {
	go func() { m.buf <- buffer{msg, s} }()
}

func (m *message) run() {
	for {
		buf := <-m.buf
		x, y := m.LeftTop()
		display := fmt.Sprintf("[%d] %s", len(m.buf)+1, buf.string)
		display = runewidth.Truncate(display, m.Width(), "")
		widget.SetCells(x, y, display, buf.Style)

		m.display = true
		widget.Show()
		<-time.NewTimer(m.sec * time.Second).C
		m.Clear()
		widget.Show()
		m.display = false
	}
}
