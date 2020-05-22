// +build !js,!appengine

package runewidth

import (
	"sort"
	"testing"
	"unicode/utf8"
)

var _ sort.Interface = (*table)(nil) // ensure that type "table" does implement sort.Interface

func (t table) Len() int {
	return len(t)
}

func (t table) Less(i, j int) bool {
	return t[i].first < t[j].first
}

func (t *table) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

type tableInfo struct {
	tbl  table
	name string
}

func checkInterval(first, last rune) bool {
	return first >= 0 && first <= utf8.MaxRune &&
		last >= 0 && last <= utf8.MaxRune &&
		first <= last
}

func isCompact(t *testing.T, ti *tableInfo) bool {
	tbl := ti.tbl
	for i := range tbl {
		e := tbl[i]
		if !checkInterval(e.first, e.last) { // sanity check
			t.Errorf("table invalid: table = %s index = %d %v", ti.name, i, e)
			return false
		}
		if i+1 < len(tbl) && e.last+1 >= tbl[i+1].first { // can be combined into one entry
			t.Errorf("table not compact: table = %s index = %d %v %v", ti.name, i, e, tbl[i+1])
			return false
		}
	}
	return true
}

var runewidthtests = []struct {
	in  rune
	out int
}{
	{'世', 2},
	{'界', 2},
	{'ｾ', 1},
	{'ｶ', 1},
	{'ｲ', 1},
	{'☆', 2}, // double width in ambiguous
	{'☺', 1},
	{'☻', 1},
	{'♥', 2},
	{'♦', 1},
	{'♣', 2},
	{'♠', 2},
	{'♂', 2},
	{'♀', 2},
	{'♪', 2},
	{'♫', 1},
	{'☼', 1},
	{'↕', 2},
	{'‼', 1},
	{'↔', 2},
	{'\x00', 0},
	{'\x01', 0},
	{'\u0300', 0},
	{'\u2028', 0},
	{'\u2029', 0},
	{'a', 1}, // ASCII classified as "na" (narrow)
	{'⟦', 1}, // non-ASCII classified as "na" (narrow)
	{'👁', 2},
}

func TestRuneWidth(t *testing.T) {
	c := NewCondition()
	for _, tt := range runewidthtests {
		if out := c.RuneWidth(tt.in); out != tt.out {
			t.Errorf("RuneWidth(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
}

var stringwidthtests = []struct {
	in  string
	out int
}{
	{"■㈱の世界①", 11},
	{"スター☆", 8},
	{"つのだ☆HIRO", 12},
}

func TestStringWidth(t *testing.T) {
	c := NewCondition()
	for _, tt := range stringwidthtests {
		if out := c.StringWidth(tt.in); out != tt.out {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, out, tt.out)
		}
	}
}

func TestStringWidthInvalid(t *testing.T) {
	s := "こんにちわ\x00世界"
	if out := StringWidth(s); out != 14 {
		t.Errorf("StringWidth(%q) = %d, want %d", s, out, 14)
	}
}

func TestTruncateSmaller(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := Truncate(s, 10, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

func TestTruncate(t *testing.T) {
	s := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."
	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 79 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 79, width)
	}
}

func TestTruncateFit(t *testing.T) {
	s := "aあいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "aあいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."

	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 80 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 80, width)
	}
}

func TestTruncateJustFit(t *testing.T) {
	s := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "あいうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"

	out := Truncate(s, 80, "...")
	if out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
	width := StringWidth(out)
	if width != 80 {
		t.Errorf("width of Truncate(%q) should be %d, but %d", s, 80, width)
	}
}

func TestWrap(t *testing.T) {
	s := `東京特許許可局局長はよく柿喰う客だ/東京特許許可局局長はよく柿喰う客だ
123456789012345678901234567890

END`
	expected := `東京特許許可局局長はよく柿喰う
客だ/東京特許許可局局長はよく
柿喰う客だ
123456789012345678901234567890

END`

	if out := Wrap(s, 30); out != expected {
		t.Errorf("Wrap(%q) = %q, want %q", s, out, expected)
	}
}

func TestTruncateNoNeeded(t *testing.T) {
	s := "あいうえおあい"
	expected := "あいうえおあい"

	if out := Truncate(s, 80, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillLeft(t *testing.T) {
	s := "あxいうえお"
	expected := "    あxいうえお"

	if out := FillLeft(s, 15); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillLeftFit(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := FillLeft(s, 10); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRight(t *testing.T) {
	s := "あxいうえお"
	expected := "あxいうえお    "

	if out := FillRight(s, 15); out != expected {
		t.Errorf("FillRight(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRightFit(t *testing.T) {
	s := "あいうえお"
	expected := "あいうえお"

	if out := FillRight(s, 10); out != expected {
		t.Errorf("FillRight(%q) = %q, want %q", s, out, expected)
	}
}

func TestEnv(t *testing.T) {
	if w := RuneWidth('│'); w != 1 {
		t.Errorf("RuneWidth('│') = %d, want %d", w, 1)
	}
}

func TestZeroWidthJoiner(t *testing.T) {
	c := NewCondition()

	var tests = []struct {
		in   string
		want int
	}{
		{"👩", 2},
		{"👩‍", 2},
		{"👩‍🍳", 2},
		{"‍🍳", 2},
		{"👨‍👨", 2},
		{"👨‍👨‍👧", 2},
		{"🏳️‍🌈", 2},
		{"あ👩‍🍳い", 6},
		{"あ‍🍳い", 6},
		{"あ‍い", 4},
	}

	for _, tt := range tests {
		if got := c.StringWidth(tt.in); got != tt.want {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
