package runewidth

import (
	"strings"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

//go:generate make

var (
	// DefaultCondition is a condition in current locale
	DefaultCondition = &Condition{}
)

func wcwidth9_lookup_width(r rune) int {
	c := int(r)

	t1_offset := wcwidth9_tables_0[c>>13&0xFF]

	t2_offset := wcwidth9_tables_1[128*t1_offset+(c>>6&0x7F)]

	packed_widths := wcwidth9_tables_2[16*t2_offset+(c>>2&0xF)]

	width := packed_widths >> (2 * (c & 0b11)) & 0b11

	if width == 3 {
		return -1
	}

	return width
}

func wcwidth9_width(r rune) int {
	if r < 0x7F {
		if r >= 0x20 {
			return 1
		}
		if r == 0 {
			return 0
		}
		return -1
	}

	if r >= 0xA0 {
		return wcwidth9_lookup_width(r)
	}

	return -1
}

type Condition struct {
	EastAsianWidth bool
}

// NewCondition return new instance of Condition which is current locale.
func NewCondition() *Condition {
	return &Condition{
		EastAsianWidth: false,
	}
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func (c *Condition) RuneWidth(r rune) int {
	if width := wcwidth9_width(r); width < 0 {
		return 0
	} else {
		return width
	}
}

// StringWidth return width as you can see
func (c *Condition) StringWidth(s string) (width int) {
	g := graphemes.FromString(s)
	for g.Next() {
		chWidth := 0
		for _, r := range g.Value() {
			chWidth += c.RuneWidth(r)
		}
		width += chWidth
	}
	return
}

// Truncate return string truncated with w cells
func (c *Condition) Truncate(s string, w int, tail string) string {
	if c.StringWidth(s) <= w {
		return s
	}
	w -= c.StringWidth(tail)
	var width int
	pos := len(s)
	g := graphemes.FromString(s)
	for g.Next() {
		chWidth := 0
		for _, r := range g.Value() {
			chWidth += c.RuneWidth(r)
		}
		if width+chWidth > w {
			pos = g.Start()
			break
		}
		width += chWidth
	}
	return s[:pos] + tail
}

// TruncateLeft cuts w cells from the beginning of the `s`.
func (c *Condition) TruncateLeft(s string, w int, prefix string) string {
	if c.StringWidth(s) <= w {
		return prefix
	}

	var width int
	pos := len(s)

	g := graphemes.FromString(s)
	for g.Next() {
		chWidth := 0
		for _, r := range g.Value() {
			chWidth += c.RuneWidth(r)
		}

		if width+chWidth > w {
			if width < w {
				pos = g.End()
				prefix += strings.Repeat(" ", width+chWidth-w)
			} else {
				pos = g.Start()
			}

			break
		}

		width += chWidth
	}

	return prefix + s[pos:]
}

// Wrap return string wrapped with w cells
func (c *Condition) Wrap(s string, w int) string {
	width := 0
	out := ""
	for _, r := range s {
		cw := c.RuneWidth(r)
		if r == '\n' {
			out += string(r)
			width = 0
			continue
		} else if width+cw > w {
			out += "\n"
			width = 0
			out += string(r)
			width += cw
			continue
		}
		out += string(r)
		width += cw
	}
	return out
}

// FillLeft return string filled in left by spaces in w cells
func (c *Condition) FillLeft(s string, w int) string {
	width := c.StringWidth(s)
	count := w - width
	if count > 0 {
		b := make([]byte, count)
		for i := range b {
			b[i] = ' '
		}
		return string(b) + s
	}
	return s
}

// FillRight return string filled in left by spaces in w cells
func (c *Condition) FillRight(s string, w int) string {
	width := c.StringWidth(s)
	count := w - width
	if count > 0 {
		b := make([]byte, count)
		for i := range b {
			b[i] = ' '
		}
		return s + string(b)
	}
	return s
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func RuneWidth(r rune) int {
	return DefaultCondition.RuneWidth(r)
}

// StringWidth return width as you can see
func StringWidth(s string) (width int) {
	return DefaultCondition.StringWidth(s)
}

// Truncate return string truncated with w cells
func Truncate(s string, w int, tail string) string {
	return DefaultCondition.Truncate(s, w, tail)
}

// TruncateLeft cuts w cells from the beginning of the `s`.
func TruncateLeft(s string, w int, prefix string) string {
	return DefaultCondition.TruncateLeft(s, w, prefix)
}

// Wrap return string wrapped with w cells
func Wrap(s string, w int) string {
	return DefaultCondition.Wrap(s, w)
}

// FillLeft return string filled in left by spaces in w cells
func FillLeft(s string, w int) string {
	return DefaultCondition.FillLeft(s, w)
}

// FillRight return string filled in left by spaces in w cells
func FillRight(s string, w int) string {
	return DefaultCondition.FillRight(s, w)
}
