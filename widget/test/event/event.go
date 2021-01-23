package main

import (
	"fmt"

	"github.com/anmitsu/goful/widget"
	"github.com/gdamore/tcell/v2"
)

func main() {
	widget.Init()
	defer widget.Fini()

	fmt.Print("Exit by q; ")
	for {
		switch ev := widget.PollEvent().(type) {
		case *tcell.EventKey:
			key := widget.EventToString(ev)
			if key == "q" {
				return
			}
			fmt.Printf("key %d rune %c name %s -> %s; ",
				ev.Key(), ev.Rune(), ev.Name(), key)
		}
	}
}
