// Package goful is CUI file manager.
package goful

import (
	"github.com/anmitsu/goful/filer"
	"github.com/anmitsu/goful/infobar"
	"github.com/anmitsu/goful/menu"
	"github.com/anmitsu/goful/message"
	"github.com/anmitsu/goful/progbar"
	"github.com/anmitsu/goful/widget"
	"github.com/nsf/termbox-go"
)

// Goful is CUI file manager.
type Goful struct {
	*filer.Filer
	next      widget.Widget
	shell     shell
	event     chan termbox.Event
	interrupt chan int
	callback  chan func()
	exit      bool
}

// New creates a new goful client based on file recording the previous state.
func New(file string) *Goful {
	width, height := termbox.Size()
	goful := &Goful{
		Filer:     filer.NewFromJSON(file, 0, 0, width, height-2),
		next:      nil,
		shell:     shell{"", []string{}, "", []string{}},
		event:     make(chan termbox.Event, 20),
		interrupt: make(chan int, 2),
		callback:  make(chan func()),
		exit:      false,
	}
	goful.Workspace().ReloadAll()
	return goful
}

// ConfigFiler sets keymap function for filer.
func (g *Goful) ConfigFiler(f func(*Goful) widget.Keymap) {
	g.MergeKeymap(f(g))
}

// ConfigShell sets shell name and options.
func (g *Goful) ConfigShell(shell string, opts ...string) {
	g.shell.name = shell
	g.shell.opts = opts
}

// ConfigScreen sets screen name and options.
func (g *Goful) ConfigScreen(screen string, opts ...string) {
	g.shell.screen = screen
	g.shell.screenopts = opts
}

// Next returns next widget for drawing and input.
func (g *Goful) Next() widget.Widget { return g.next }

// Disconnect references to next widget for exiting.
func (g *Goful) Disconnect() { g.next = nil }

// Resize implements widget.Widget and resizes all widget.
func (g *Goful) Resize(x, y, width, height int) {
	g.Filer.Resize(x, y, width, height-2)
	if wid := g.Next(); wid != nil {
		wid.Resize(x, y, width, height-2)
	}
	message.Resize(width, height)
	infobar.Resize(width, height)
	progbar.Resize(width, height)
}

// Draw implements widget.Widget and draws all widget.
func (g *Goful) Draw() {
	g.Filer.Draw()
	if g.Next() != nil {
		g.Next().Draw()
	}
	message.Draw()
	infobar.Draw(g.File())
	progbar.Draw()
}

// Input implements widget.Widget.
func (g *Goful) Input(key string) {
	if g.Next() != nil {
		g.Next().Input(key)
	} else {
		g.Filer.Input(key)
	}
}

// Menu runs menu mode.
func (g *Goful) Menu(name string) {
	m, err := menu.New(name, g)
	if err != nil {
		message.Error(err)
		return
	}
	g.next = m
}

// Run goful client.
func (g *Goful) Run() {
	message.Info("Welcome to goful")

	go func() {
		for {
			g.event <- termbox.PollEvent()
		}
	}()

	g.Draw()
	widget.Flush()

	for !g.exit {
		select {
		case ev := <-g.event:
			g.eventHandler(ev)
		case <-g.interrupt:
			<-g.interrupt
		case callback := <-g.callback:
			callback()
		}
		g.Draw()
		widget.Flush()
	}
}

func (g *Goful) quit() {
	g.exit = true
}

func (g *Goful) syncCallback(callback func()) {
	g.callback <- callback
}

func (g *Goful) eventHandler(ev termbox.Event) {
	if key := widget.Event2String(&ev); key != "" {
		if key == "resize" {
			width, height := ev.Width, ev.Height
			g.Resize(0, 0, width, height)
		} else {
			g.Input(key)
		}
	}
}

func (g *Goful) dialog(message string, options ...string) string {
	g.interrupt <- 1
	defer func() { g.interrupt <- 1 }()

	tmp := g.Next()
	dialog := g.DialogMode(message, options...)
	g.Draw()
	widget.Flush()

	for g.Next() != nil {
		select {
		case ev := <-g.event:
			g.eventHandler(ev)
		}
		g.Draw()
		widget.Flush()
	}
	g.next = tmp
	return dialog.Result()
}
