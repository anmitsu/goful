package message

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/utils"
	"github.com/anmitsu/goful/widget"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Info notifies an infomation message.
func Info(s string) {
	messenger.info(s)
}

// Infof notifies a formated infomation message.
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
func Resize(width, height int) {
	messenger.Resize(0, height-2, width, 1)
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

// SetInfoLog sets a path to log infomation messages.
// If sets to "", non logging.
func SetInfoLog(path string) {
	messenger.infolog = utils.ExpandPath(path)
}

// SetErrorLog sets a path to log error messages.
// If sets to "", non logging.
func SetErrorLog(path string) {
	messenger.errlog = utils.ExpandPath(path)
}

// Init initializes the message works asynchronously.
func Init() {
	width, height := termbox.Size()
	messenger = &message{
		Window:  widget.NewWindow(0, height-2, width, 1),
		buf:     make(chan buffer, 20),
		display: false,
		sec:     5,
	}
	go messenger.run()
}

var messenger *message

type buffer struct {
	string
	look.Look
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
			m.add(err.Error(), look.MessageErr())
		}
	}
	m.add(s, look.MessageInfo())
}

func (m *message) error(e error) {
	if m.errlog != "" {
		if err := m.log(m.errlog, e.Error()); err != nil {
			m.add(err.Error(), look.MessageErr())
		}
	}
	m.add(e.Error(), look.MessageErr())
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

func (m *message) add(msg string, l look.Look) {
	m.buf <- buffer{msg, l}
}

func (m *message) run() {
	for {
		buf := <-m.buf
		x, y := m.LeftTop()
		x++
		display := fmt.Sprintf("[%d] %s", len(m.buf)+1, buf.string)
		display = runewidth.Truncate(display, m.Width(), "")
		widget.SetCells(x, y, display, buf.Look)

		m.display = true
		widget.Flush()
		<-time.NewTimer(m.sec * time.Second).C
		m.Clear()
		widget.Flush()
		m.display = false
	}
}
