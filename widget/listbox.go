package widget

import (
	"fmt"

	"github.com/anmitsu/goful/look"
	"github.com/anmitsu/goful/util"
	"github.com/mattn/go-runewidth"
)

// Drawer describes drawring for list box contents.
type Drawer interface {
	Name() string
	Draw(x, y, width int, focus bool)
}

type content struct {
	name string
}

func newContent(name string) *content {
	return &content{name}
}

func (e *content) Name() string { return e.name }

func (e *content) Draw(x, y, width int, focus bool) {
	s := runewidth.Truncate(e.Name(), width, "~")
	s = runewidth.FillRight(s, width)
	style := look.Default()
	if focus {
		style = style.Reverse(true)
	}
	SetCells(x, y, s, style)
}

type highlightContent struct {
	name      string
	highlight string
}

func newHighlightContent(name, highlight string) *highlightContent {
	return &highlightContent{name, highlight}
}

func (e *highlightContent) Name() string { return e.name }

func (e *highlightContent) Draw(x, y, width int, focus bool) {
	style := look.Default()
	if focus {
		style = style.Reverse(true)
	}
	if e.highlight == "" {
		s := runewidth.Truncate(e.Name(), width, "~")
		s = runewidth.FillRight(s, width)
		SetCells(x, y, s, style)
		return
	}

	name := runewidth.Truncate(e.Name(), width, "~")
	for _, s := range util.SplitWithSep(name, e.highlight) {
		if s == e.highlight {
			style := look.Highlight()
			if focus {
				style = style.Reverse(true)
			}
			x = SetCells(x, y, s, style)
		} else {
			x = SetCells(x, y, s, style)
		}
	}
	if w := runewidth.StringWidth(e.Name()); w < width {
		s := runewidth.FillRight("", width-w)
		SetCells(x, y, s, style)
	}
}

// ListBox is a scrollable window listing contents.
type ListBox struct {
	*Window
	list   []Drawer
	cursor int // index of a current content in the list
	offset int // index of a top position to display
	title  string
	lower  int
	column int
}

// NewListBox creates a new list box specified coordinates and sizes.
func NewListBox(x, y, width, height int, title string) *ListBox {
	return &ListBox{
		Window: NewWindow(x, y, width, height),
		list:   []Drawer{},
		cursor: 0,
		offset: 0,
		title:  title,
		lower:  0,
		column: 1,
	}
}

func (b *ListBox) Len() int {
	return len(b.list)
}

func (b *ListBox) Swap(i, j int) {
	b.list[i], b.list[j] = b.list[j], b.list[i]
}

func (b *ListBox) Less(i, j int) bool {
	return b.list[i].Name() < b.list[j].Name()
}

// List returns contents in the list box.
func (b *ListBox) List() []Drawer { return b.list }

// SetList sets contents to the list box.
func (b *ListBox) SetList(l []Drawer) { b.list = l }

// ClearList empties the list.
func (b *ListBox) ClearList() { b.list = []Drawer{} }

// IsEmpty reports whether the list empty.
func (b *ListBox) IsEmpty() bool { return len(b.list) == 0 }

// AppendList appends new contents to the end.
func (b *ListBox) AppendList(e ...Drawer) {
	b.list = append(b.list, e...)
}

// AppendString appends new contents by strings to the end.
func (b *ListBox) AppendString(s ...string) {
	for _, v := range s {
		b.list = append(b.list, newContent(v))
	}
}

// AppendHighlightString appends a new highlighted content by the string to the end.
func (b *ListBox) AppendHighlightString(s, highlight string) {
	b.list = append(b.list, newHighlightContent(s, highlight))
}

// CurrentContent gets the content on the cursor.
func (b *ListBox) CurrentContent() Drawer {
	return b.list[b.cursor]
}

// Title returns the list box title.
func (b *ListBox) Title() string { return b.title }

// SetTitle sets the list box title.
func (b *ListBox) SetTitle(s string) { b.title = s }

// Lower returns a lower bound that the cursor can move.
func (b *ListBox) Lower() int { return b.lower }

// SetLower sets the lower.
func (b *ListBox) SetLower(lb int) { b.lower = lb }

// Upper returns an upper bound that the cursor can move.
func (b *ListBox) Upper() int { return len(b.list) }

// Column returns the column of the list box.
func (b *ListBox) Column() int { return b.column }

// SetColumn sets the column.
func (b *ListBox) SetColumn(n int) { b.column = n }

// ColumnAdjustContentsWidth adjusts the column by list box constent widths.
// Fit to the widest content.
func (b *ListBox) ColumnAdjustContentsWidth() {
	column := b.Width()
	for _, e := range b.list {
		w := runewidth.StringWidth(e.Name())
		if n := b.Width() / (w + 1); n < column && n > 0 {
			column = n
		}
	}
	b.column = column
}

// Offset returns a top position to display.
func (b *ListBox) Offset() int { return b.offset }

func (b *ListBox) isScroll() bool {
	return b.Upper() > (b.Height()-2)*b.column
}

// SetOffsetCenteredCursor sets the offset so that centred cursor.
func (b *ListBox) SetOffsetCenteredCursor() {
	if !b.isScroll() {
		return
	}
	b.offset = b.cursor - b.rowCol()/2

	if b.offset < 0 {
		b.offset = 0
	} else if b.offset > b.Upper() {
		b.offset = b.Upper()
	}
}

// Cursor returns the current position of contents.
func (b *ListBox) Cursor() int { return b.cursor }

// SetCursor sets the cursor.
func (b *ListBox) SetCursor(x int) {
	if x < b.Lower() {
		b.cursor = b.Lower()
	} else if x > b.Upper()-1 {
		b.cursor = b.Upper() - 1
	} else {
		b.cursor = x
	}
}

// MoveCursor moves the cursor.
func (b *ListBox) MoveCursor(amount int) {
	b.cursor += amount
	if b.cursor < b.Lower() {
		b.cursor = b.Lower()
	} else if b.cursor > b.Upper()-1 {
		b.cursor = b.Upper() - 1
	}
}

// CursorDown moves the cursor to the down.
func (b *ListBox) CursorDown() {
	if b.Upper()-b.column-1 >= b.cursor {
		b.MoveCursor(b.column)
	}
}

// CursorUp moves the cursor to the up.
func (b *ListBox) CursorUp() {
	if b.Lower()+b.column <= b.cursor {
		b.MoveCursor(-b.column)
	}
}

// CursorToRight moves the cursor to the right.
func (b *ListBox) CursorToRight() {
	if b.column > 1 {
		b.MoveCursor(1)
	}
}

// CursorToLeft moves the cursor to the left.
func (b *ListBox) CursorToLeft() {
	if b.column > 1 {
		b.MoveCursor(-1)
	}
}

// SetCursorByName sets the cursor by the content name.
func (b *ListBox) SetCursorByName(name string) {
	if idx := b.IndexByName(name); idx != -1 {
		b.SetCursor(idx)
	}
}

// IndexByName returns the index to match content name.
func (b *ListBox) IndexByName(name string) int {
	for i, content := range b.list {
		if name == content.Name() {
			return i
		}
	}
	return b.lower
}

// MoveTop moves cursor to index of lower list.
func (b *ListBox) MoveTop() {
	b.SetCursor(b.Lower())
}

// MoveBottom moves cursor to index of upper list.
func (b *ListBox) MoveBottom() {
	b.SetCursor(b.Upper() - 1)
}

// PageDown moves cursor to next page.
func (b *ListBox) PageDown() {
	height := b.rowCol()
	if b.offset+height >= b.Upper() {
		return
	}
	b.offset += height
	b.cursor += height
}

// PageUp moves cursor to previous page.
func (b *ListBox) PageUp() {
	if b.offset == 0 {
		return
	}
	height := b.rowCol()
	if b.offset-height < 0 {
		b.offset = 0
		b.cursor = 0
	} else {
		b.offset -= height
		b.cursor -= height
	}
}

// Scroll scrolls the list box.
func (b *ListBox) Scroll(amount int) {
	amount *= b.column
	bottom := b.offset + b.rowCol()

	if amount > 0 {
		if bottom >= b.Upper() {
			return
		}
		b.offset += amount
		if b.cursor < b.offset {
			b.cursor += amount
		}
	} else {
		if b.offset != 0 {
			b.offset += amount
			bottom += amount
			if b.cursor >= bottom {
				b.cursor += amount
			}
		}
	}
}

// AdjustCursor adjusts the cursor within range of upper and lower.
func (b *ListBox) AdjustCursor() {
	if b.cursor >= b.Upper() {
		b.cursor = b.Upper() - 1
	} else if b.cursor < b.Lower() {
		b.cursor = b.Lower()
	}
}

func (b *ListBox) rowCol() int {
	return (b.Height() - 2) * b.column
}

// AdjustOffset adjusts the offset within range of upper and lower.
func (b *ListBox) AdjustOffset() {
	if !b.isScroll() {
		b.offset = 0
		return
	}
	if b.rowCol()/2 < 1 {
		return
	}

	for b.cursor < b.offset {
		b.offset -= b.rowCol() / 2
	}
	for b.cursor >= b.offset+b.rowCol() {
		b.offset += b.rowCol() / 2
	}

	if b.offset < 0 {
		b.offset = 0
	} else if b.offset > b.Upper() {
		b.offset = b.Upper()
	}
}

// ScrollRate returns rate of offset.
func (b *ListBox) ScrollRate() string {
	base := float64(b.Upper() - b.rowCol())
	if base == 0 {
		base = float64(b.Upper())
	}
	p := float64(b.offset) / base * 100
	if p == 0 {
		return "Top"
	} else if p >= 100 {
		return "Bot"
	} else {
		return fmt.Sprintf("%d%s", int(p), "%")
	}
}

func (b *ListBox) drawHeader() {
	title := fmt.Sprintf("%s [%d/%d] %s", b.title, b.cursor+1, b.Upper(), b.ScrollRate())
	x, y := b.LeftTop()
	SetCells(x, y, title, look.Title())
}

func (b *ListBox) drawScrollbar() {
	if !b.isScroll() {
		return
	}
	height := b.Height() - 2
	offset := int(float64(b.offset) / float64(b.Upper()-b.rowCol()) * float64(height))
	if offset > height-1 {
		offset = height - 1
	}

	x, y := b.RightTop()
	y++
	for i := 0; i < height; i++ {
		if i == offset {
			SetCells(x, y+i, "=", look.Default())
		} else {
			SetCells(x, y+i, "|", look.Default())
		}
	}
}

// Draw list box with title and scrollbar.
func (b *ListBox) Draw() {
	if b.Upper() < 1 {
		return
	}
	b.AdjustCursor()
	b.AdjustOffset()
	b.Clear()
	b.Border()
	b.drawHeader()
	b.drawScrollbar()

	width, height := b.Width()-2, b.Height()-2
	shift := 1
	if b.border == AllBorder {
		shift++
	}
	colwidth := width/b.column - shift + 1
	row, col := 1, 0
	for i := b.offset; i < b.Upper(); i++ {
		if col >= b.column {
			col = 0
			row++
			if row > height {
				break
			}
		}
		x, y := b.LeftTop()
		x += col*colwidth + shift
		y += row
		if i != b.cursor {
			b.list[i].Draw(x, y, colwidth, false)
		} else {
			b.list[i].Draw(x, y, colwidth, true)
		}
		col++
	}
}
