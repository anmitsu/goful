// Package menu provides the menu window widget.
package menu

import (
	"fmt"

	"github.com/anmitsu/goful/widget"
)

var menusMap = map[string][]*menuItem{}

// Add menu items as label, acceleration key and callback function and
// the number of arguments `a' must be a multiple of three.
func Add(name string, a ...interface{}) {
	if len(a)%3 != 0 {
		panic("items must be a multiple of three")
	}
	items := menusMap[name]
	for i := 0; i < len(a); i += 3 {
		accel := a[i].(string)
		label := a[i+1].(string)
		callback := a[i+2].(func())
		items = append(items, &menuItem{accel, label, callback})
	}
	menusMap[name] = items
}

var keymap func(*Menu) widget.Keymap

// Config the keymap function for a menu.
func Config(config func(*Menu) widget.Keymap) {
	keymap = config
}

type menuItem struct {
	accel    string
	label    string
	callback func()
}

// Menu is a list box to execute for a acceleration key.
type Menu struct {
	*widget.ListBox
	filer widget.Widget
}

// New creates a new menu based on filer widget sizes.
func New(name string, filer widget.Widget) (*Menu, error) {
	items, ok := menusMap[name]
	if !ok {
		return nil, fmt.Errorf("Not found menu `%s'", name)
	}
	x, y := filer.LeftBottom()
	width := filer.Width()
	height := len(items) + 2
	if max := filer.Height() / 2; height > max {
		height = max
	}
	menu := &Menu{
		ListBox: widget.NewListBox(x, y-height+1, width, height, name),
		filer:   filer,
	}
	for _, item := range items {
		s := fmt.Sprintf("%-3s %s", item.accel, item.label)
		menu.AppendString(s)
	}
	return menu, nil
}

// Resize the menu window.
func (w *Menu) Resize(x, y, width, height int) {
	h := len(menusMap[w.Title()]) + 2
	if max := height / 2; h > max {
		h = max
	}
	w.ListBox.Resize(x, height-h, width, h)
}

// Exec executes a menu item on the cursor and exits the menu.
func (w *Menu) Exec() {
	w.Exit()
	menusMap[w.Title()][w.Cursor()].callback()
}

// Input to the list box or execute a menu item with the acceleration key.
func (w *Menu) Input(key string) {
	keymap := keymap(w)
	if callback, ok := keymap[key]; ok {
		callback()
	} else {
		for _, item := range menusMap[w.Title()] {
			if item.accel == key {
				w.Exit()
				item.callback()
			}
		}
	}
}

// Exit the menu mode.
func (w *Menu) Exit() { w.filer.Disconnect() }

// Next implements widget.Widget.
func (w *Menu) Next() widget.Widget { return widget.Nil() }

// Disconnect implements widget.Widget.
func (w *Menu) Disconnect() {}
