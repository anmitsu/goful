package filer

import (
	"log"
	"os"

	"github.com/anmitsu/goful/widget"
)

type layoutType int

const (
	layoutTile layoutType = iota
	layoutTileTop
	layoutTileBottom
	layoutOneline
	layoutOneColumn
	layoutFullscreen
)

// Workspace is a box storing and layouting directories.
type Workspace struct {
	*widget.Window
	Dirs   []*Directory `json:"directories"`
	Layout layoutType   `json:"layout"`
	Title  string       `json:"title"`
	Focus  int          `json:"focus"`
}

// NewWorkspace returns a new workspace of specified sizes.
func NewWorkspace(x, y, width, height int, title string) *Workspace {
	return &Workspace{
		widget.NewWindow(x, y, width, height),
		[]*Directory{},
		layoutTile,
		title,
		0,
	}
}

func (w *Workspace) init4json(x, y, width, height int) {
	w.Window = widget.NewWindow(x, y, width, height)
}

// CreateDir adds the home directory to the head.
func (w *Workspace) CreateDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := NewDirectory(0, 0, 0, 0)
	dir.Chdir(home)
	w.Dirs = append(w.Dirs, nil)
	copy(w.Dirs[1:], w.Dirs[:len(w.Dirs)-1])
	w.Dirs[0] = dir
	w.SetFocus(0)
	w.allocate()
}

// CloseDir closes the focused directory.
func (w *Workspace) CloseDir() {
	if len(w.Dirs) < 2 {
		return
	}
	i := w.Focus
	w.Dirs = append(w.Dirs[:i], w.Dirs[i+1:]...)
	if w.Focus >= len(w.Dirs) {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
	w.allocate()
}

// ChdirNeighbor changes the focused path a neighbor directory path.
func (w *Workspace) ChdirNeighbor() {
	w.Dir().Chdir(w.NextDir().Path)
}

func (w *Workspace) visible(visible bool) {
	if visible {
		w.ReloadAll()
	} else {
		for _, d := range w.Dirs {
			d.ClearList()
		}
	}
}

// MoveFocus moves the focus with specified amounts.
func (w *Workspace) MoveFocus(amount int) {
	w.Focus += amount
	if len(w.Dirs) <= w.Focus {
		w.Focus = 0
	} else if w.Focus < 0 {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
}

// SetFocus sets the focus to a specified position.
func (w *Workspace) SetFocus(x int) {
	w.Focus = x
	if w.Focus < 0 {
		w.Focus = 0
	} else if w.Focus > len(w.Dirs)-1 {
		w.Focus = len(w.Dirs) - 1
	}
	w.attach()
}

func (w *Workspace) attach() {
	err := os.Chdir(w.Dir().Path)
	if err != nil {
		log.Fatalln(err)
	}
}

// ReloadAll reloads all directories.
func (w *Workspace) ReloadAll() {
	for _, d := range w.Dirs {
		d.reload()
	}
	err := os.Chdir(w.Dir().Path)
	if err != nil {
		log.Fatalln(err)
	}
}

// Dir returns the focused directory.
func (w *Workspace) Dir() *Directory {
	return w.Dirs[w.Focus]
}

// NextDir returns the next directory.
func (w *Workspace) NextDir() *Directory {
	return w.Dirs[w.nextIndex()]
}

// PrevDir returns the previous directory.
func (w *Workspace) PrevDir() *Directory {
	return w.Dirs[w.prevIndex()]
}

// SwapNextDir swaps focus and next directories.
func (w *Workspace) SwapNextDir() {
	next := w.nextIndex()
	w.Dirs[w.Focus], w.Dirs[next] = w.Dirs[next], w.Dirs[w.Focus]
	w.MoveFocus(1)
	w.allocate()
}

// SwapPrevDir swaps focus and previous directories.
func (w *Workspace) SwapPrevDir() {
	prev := w.prevIndex()
	w.Dirs[w.Focus], w.Dirs[prev] = w.Dirs[prev], w.Dirs[w.Focus]
	w.MoveFocus(-1)
	w.allocate()
}

func (w *Workspace) nextIndex() int {
	i := w.Focus + 1
	if i >= len(w.Dirs) {
		return 0
	}
	return i
}

func (w *Workspace) prevIndex() int {
	i := w.Focus - 1
	if i < 0 {
		return len(w.Dirs) - 1
	}
	return i
}

// SetTitle sets the workspace title.
func (w *Workspace) SetTitle(title string) {
	w.Title = title
}

// LayoutTile allocates to the tile layout.
func (w *Workspace) LayoutTile() {
	w.Layout = layoutTile
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	width := w.Width() / 2
	w.Dirs[0].Resize(x, y, width, w.Height())
	height := w.Height() / k
	hodd := w.Height() % k
	wodd := w.Width() % 2
	for i, d := range w.Dirs[1:k] {
		d.Resize(x+width, y+height*i, width+wodd, height)
	}
	w.Dirs[k].Resize(x+width, y+height*(k-1), width+wodd, height+hodd)
}

// LayoutTileTop allocates to the tile top layout.
func (w *Workspace) LayoutTileTop() {
	w.Layout = layoutTileTop
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	height := w.Height() / 2
	hodd := w.Height() % 2

	width := w.Width() / k
	wodd := w.Width() % 2

	w.Dirs[0].Resize(x, y, width, height)
	w.Dirs[k].Resize(x, y+height, w.Width(), height+hodd)
	if k < 2 {
		return
	}
	for i, d := range w.Dirs[1 : k-1] {
		d.Resize(x+width*(i+1), y, width, height)
	}
	w.Dirs[k-1].Resize(x+width*(k-1), y, width+wodd, height)
}

// LayoutTileBottom allocates to the tile bottom layout.
func (w *Workspace) LayoutTileBottom() {
	w.Layout = layoutTileBottom
	x, y := w.LeftTop()
	k := len(w.Dirs) - 1
	if k < 1 {
		w.Dirs[0].Resize(x, y, w.Width(), w.Height())
		return
	}
	height := w.Height() / 2
	hodd := w.Height() % 2

	w.Dirs[0].Resize(x, y, w.Width(), height)

	width := w.Width() / k
	for i, d := range w.Dirs[1:k] {
		d.Resize(x+width*i, y+height, width, height+hodd)
	}
	wodd := w.Width() % 2
	w.Dirs[k].Resize(x+width*(k-1), y+height, width+wodd, height+hodd)
}

// LayoutOnerow allocates to the one line layout.
func (w *Workspace) LayoutOnerow() {
	w.Layout = layoutOneline
	x, y := w.LeftTop()
	k := len(w.Dirs)
	width := w.Width() / k
	for i, d := range w.Dirs[:k-1] {
		d.Resize(x+width*i, y, width, w.Height())
	}
	wodd := w.Width() % k
	w.Dirs[k-1].Resize(x+width*(k-1), y, width+wodd, w.Height())
}

// LayoutOnecolumn allocates to the one column layout.
func (w *Workspace) LayoutOnecolumn() {
	w.Layout = layoutOneColumn
	x, y := w.LeftTop()
	k := len(w.Dirs)
	height := w.Height() / k
	for i, d := range w.Dirs[:k-1] {
		d.Resize(x, y+height*i, w.Width(), height)
	}
	hodd := w.Height() % k
	w.Dirs[k-1].Resize(x, y+height*(k-1), w.Width(), height+hodd)
}

// LayoutFullscreen allocates to the full screen layout.
func (w *Workspace) LayoutFullscreen() {
	w.Layout = layoutFullscreen
	for _, d := range w.Dirs {
		x, y := w.LeftTop()
		d.Resize(x, y, w.Width(), w.Height())
	}
}

func (w *Workspace) allocate() {
	switch w.Layout {
	case layoutTile:
		w.LayoutTile()
	case layoutTileTop:
		w.LayoutTileTop()
	case layoutTileBottom:
		w.LayoutTileBottom()
	case layoutOneline:
		w.LayoutOnerow()
	case layoutOneColumn:
		w.LayoutOnecolumn()
	case layoutFullscreen:
		w.LayoutFullscreen()
	}
}

// Resize and layout allocates.
func (w *Workspace) Resize(x, y, width, height int) {
	w.Window.Resize(x, y, width, height)
	w.allocate()
}

// ResizeRelative relative resizes and layout allocates.
func (w *Workspace) ResizeRelative(x, y, width, height int) {
	w.Window.ResizeRelative(x, y, width, height)
	w.allocate()
}

// Draw all directories.
func (w *Workspace) Draw() {
	switch w.Layout {
	case layoutFullscreen:
		w.Dir().draw(true)
	default:
		for i, d := range w.Dirs {
			if i != w.Focus {
				d.draw(false)
			} else {
				d.draw(true)
			}
		}
	}
}
