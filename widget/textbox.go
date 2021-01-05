package widget

import (
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// TextBox is editable text box.
type TextBox struct {
	*Window
	text     []byte
	offset   int // editing position of []byte
	cursor   int // editing position of visualized
	Edithook func()
}

// NewTextBox returns the text box of specified size and coordinates.
func NewTextBox(x, y, width, height int) *TextBox {
	return &TextBox{
		Window:   NewWindow(x, y, width, height),
		text:     []byte{},
		offset:   0,
		cursor:   0,
		Edithook: func() {},
	}
}

// Cursor returns the editting visualized position.
func (b *TextBox) Cursor() int { return b.cursor }

// String returns the text as string.
func (b *TextBox) String() string { return string(b.text) }

// TextBeforeCursor returns the text before the cursor.
func (b *TextBox) TextBeforeCursor() string {
	return string(b.text[:b.offset])
}

// TextAfterCursor returns the text after the cursor.
func (b *TextBox) TextAfterCursor() string {
	return string(b.text[b.offset:])
}

// WidthTextBeforeCursor returns width of text before the cursor.
func (b *TextBox) WidthTextBeforeCursor() int {
	return runewidth.StringWidth(b.TextBeforeCursor())
}

// WidthTextAfterCursor returns width of text after the cursor.
func (b *TextBox) WidthTextAfterCursor() int {
	return runewidth.StringWidth(b.TextAfterCursor())
}

// MoveTop sets cursor position to the top.
func (b *TextBox) MoveTop() {
	b.offset = 0
	b.cursor = 0
}

// MoveBottom sets cursor position to the bottom.
func (b *TextBox) MoveBottom() {
	b.cursor = 0
	text := b.text[:]
	for len(text) > 0 {
		r, rsize := utf8.DecodeRune(text)
		text = text[rsize:]
		b.cursor += runewidth.RuneWidth(r)
	}
	b.offset = len(b.text)
}

// MoveCursor moves editing cursor.
func (b *TextBox) MoveCursor(offset int) {
	if offset < 0 {
		for offset < 0 {
			b.BackwardChar()
			offset++
		}
	} else if offset > 0 {
		for offset > 0 {
			b.ForwardChar()
			offset--
		}
	}
}

// ForwardChar move forwards the editing cursor only a character.
func (b *TextBox) ForwardChar() {
	text := b.text[b.offset:]
	if len(text) > 0 {
		r, rsize := utf8.DecodeRune(text)
		b.offset += rsize
		b.cursor += runewidth.RuneWidth(r)
	}
}

// BackwardChar move backwards the editing cursor only a character.
func (b *TextBox) BackwardChar() {
	text := b.text[:b.offset]
	if len(text) > 0 {
		r, rsize := utf8.DecodeLastRune(text)
		b.offset -= rsize
		b.cursor -= runewidth.RuneWidth(r)
	}
}

// ForwardWord move forwards the editing cursor only a word.
func (b *TextBox) ForwardWord() {
	text := b.text[b.offset:]

	r, rsize := utf8.DecodeRune(text)
	for !isWord(r) && len(text) > 0 {
		b.offset += rsize
		b.cursor += runewidth.RuneWidth(r)
		text = text[rsize:]
		r, rsize = utf8.DecodeRune(text)
	}

	for isWord(r) && len(text) > 0 {
		b.offset += rsize
		b.cursor += runewidth.RuneWidth(r)
		text = text[rsize:]
		r, rsize = utf8.DecodeRune(text)
	}
}

// BackwardWord move backwards the editing cursor only a word.
func (b *TextBox) BackwardWord() {
	text := b.text[:b.offset]

	r, rsize := utf8.DecodeLastRune(text)
	for !isWord(r) && len(text) > 0 {
		b.offset -= rsize
		b.cursor -= runewidth.RuneWidth(r)
		text = text[:len(text)-rsize]
		r, rsize = utf8.DecodeLastRune(text)
	}

	for isWord(r) && len(text) > 0 {
		b.offset -= rsize
		b.cursor -= runewidth.RuneWidth(r)
		text = text[:len(text)-rsize]
		r, rsize = utf8.DecodeLastRune(text)
	}
}

// DeleteChar deletes a character on editing cursor.
func (b *TextBox) DeleteChar() {
	text := b.text[b.offset:]
	if len(text) > 0 {
		_, rsize := utf8.DecodeRune(text)
		text = text[:rsize]
		b.text = DeleteBytes(b.text, b.offset, len(text))
	}
	b.Edithook()
}

// DeleteBackwardChar deletes a character on backward of editing cursor.
func (b *TextBox) DeleteBackwardChar() {
	if b.offset > 0 {
		b.BackwardChar()
		b.DeleteChar()
	}
	b.Edithook()
}

// DeleteForwardWord deletes a word on forward of editing cursor.
func (b *TextBox) DeleteForwardWord() {
	offset := b.offset
	text := b.text[offset:]

	r, rsize := utf8.DecodeRune(text)
	for !isWord(r) && len(text) > 0 {
		offset += rsize
		text = text[rsize:]
		r, rsize = utf8.DecodeRune(text)
	}

	for isWord(r) && len(text) > 0 {
		offset += rsize
		text = text[rsize:]
		r, rsize = utf8.DecodeRune(text)
	}

	text = b.text[b.offset:offset]
	b.text = DeleteBytes(b.text, b.offset, len(text))
	b.Edithook()
}

// DeleteBackwardWord deletes a word on backward of editing cursor.
func (b *TextBox) DeleteBackwardWord() {
	offset := b.offset
	text := b.text[:offset]

	r, rsize := utf8.DecodeLastRune(text)
	for !isWord(r) && len(text) > 0 {
		b.offset -= rsize
		b.cursor -= runewidth.RuneWidth(r)
		text = text[:len(text)-rsize]
		r, rsize = utf8.DecodeLastRune(text)
	}

	for isWord(r) && len(text) > 0 {
		b.offset -= rsize
		b.cursor -= runewidth.RuneWidth(r)
		text = text[:len(text)-rsize]
		r, rsize = utf8.DecodeLastRune(text)
	}

	text = b.text[b.offset:offset]
	b.text = DeleteBytes(b.text, b.offset, len(text))
	b.Edithook()
}

// KillLine deletes text on backward of editing cursor in line.
func (b *TextBox) KillLine() {
	if text := b.text[b.offset:]; len(text) > 0 {
		b.text = DeleteBytes(b.text, b.offset, len(text))
	}
	b.Edithook()
}

// KillLineAll deletes text in line.
func (b *TextBox) KillLineAll() {
	b.MoveTop()
	b.KillLine()
	b.Edithook()
}

// InsertChar inserts a character to position on editing cursor.
func (b *TextBox) InsertChar(r rune) {
	var text [utf8.UTFMax]byte
	rsize := utf8.EncodeRune(text[:], r)

	b.text = InsertBytes(b.text, text[:rsize], b.offset)

	b.offset += rsize
	b.cursor += runewidth.RuneWidth(r)
	b.Edithook()
}

// InsertString inserts string to postion on editing cursor.
func (b *TextBox) InsertString(str string) {
	for _, s := range str {
		b.InsertChar(s)
	}
}

// SetText replaces the new text.
func (b *TextBox) SetText(text string) {
	b.text = []byte(text)
	b.MoveBottom()
}

func (b *TextBox) nextChar() string {
	if len(b.text) == 0 || b.offset == len(b.text)-1 {
		return ""
	}
	return string(b.text[b.offset+1])
}

func isWord(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r)
}

func growByteSlice(s []byte, capacity int) []byte {
	if cap(s) < capacity {
		ns := make([]byte, len(s), capacity)
		copy(ns, s)
		return ns
	}
	return s
}

// InsertBytes inserts data to offset position.
func InsertBytes(s, data []byte, offset int) []byte {
	n := len(s) + len(data)
	s = growByteSlice(s, n)
	s = s[:n]
	copy(s[offset+len(data):], s[offset:])
	copy(s[offset:], data)
	return s
}

// DeleteBytes deletes bytes in offset position by length.
func DeleteBytes(s []byte, offset, length int) []byte {
	copy(s[offset:], s[offset+length:])
	s = s[:len(s)-length]
	return s
}
