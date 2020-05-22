package runewidth

import (
	"os"
)

//go:generate go run script/generate.go

var (
	// EastAsianWidth will be set true if the current locale is CJK
	EastAsianWidth bool

	// ZeroWidthJoiner is flag to set to use UTR#51 ZWJ
	ZeroWidthJoiner bool

	// DefaultCondition is a condition in current locale
	DefaultCondition = &Condition{}
)

func init() {
	handleEnv()
}

func handleEnv() {
	env := os.Getenv("RUNEWIDTH_EASTASIAN")
	if env == "" {
		EastAsianWidth = IsEastAsian()
	} else {
		EastAsianWidth = env == "1"
	}
	// update DefaultCondition
	DefaultCondition.EastAsianWidth = EastAsianWidth
	DefaultCondition.ZeroWidthJoiner = ZeroWidthJoiner
}

type interval struct {
	first rune
	last  rune
	width int
}

type table []interval

func wcwidth9_width(r rune) int {
	if r < wcwidth9_table[0].first {
		return 1
	}

	bot := 0
	top := wcwidth9_table_length
	for top >= bot {
		mid := (bot + top) >> 1

		switch {
		case wcwidth9_table[mid].last < r:
			bot = mid + 1
		case wcwidth9_table[mid].first > r:
			top = mid - 1
		default:
			return wcwidth9_table[mid].width
		}
	}

	return 1
}

// Condition have flag EastAsianWidth whether the current locale is CJK or not.
type Condition struct {
	EastAsianWidth  bool
	ZeroWidthJoiner bool
}

// NewCondition return new instance of Condition which is current locale.
func NewCondition() *Condition {
	return &Condition{
		EastAsianWidth:  EastAsianWidth,
		ZeroWidthJoiner: ZeroWidthJoiner,
	}
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func (c *Condition) RuneWidth(r rune) int {
	switch {
	case r == 0:
		return 0
	case r < 0x20:
		return 0
	case r < 0x7F:
		return 1
	case r < 0xA0:
		return 0
	case r < 0 || r > 0x10FFFF:
		return 0
	default:
		return wcwidth9_width(r)
	}
}

func (c *Condition) stringWidth(s string) (width int) {
	for _, r := range []rune(s) {
		width += c.RuneWidth(r)
	}
	return width
}

// StringWidth return width as you can see
func (c *Condition) StringWidth(s string) (width int) {
	return c.stringWidth(s)
}

// Truncate return string truncated with w cells
func (c *Condition) Truncate(s string, w int, tail string) string {
	if c.StringWidth(s) <= w {
		return s
	}
	r := []rune(s)
	tw := c.StringWidth(tail)
	w -= tw
	width := 0
	i := 0
	for ; i < len(r); i++ {
		cw := c.RuneWidth(r[i])
		if width+cw > w {
			break
		}
		width += cw
	}
	return string(r[0:i]) + tail
}

// Wrap return string wrapped with w cells
func (c *Condition) Wrap(s string, w int) string {
	width := 0
	out := ""
	for _, r := range []rune(s) {
		cw := RuneWidth(r)
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
