//go:build !js && !appengine
// +build !js,!appengine

package runewidth

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestRuneWidthChecksums(t *testing.T) {
	var testcases = []struct {
		name    string
		wantSHA string
	}{
		{"wcwidth9", "163538cabc69eea8e4f88a5a0be11a04fc335ffb8fc2e5041bf7efde4baa3179"},
	}

	for _, testcase := range testcases {
		c := NewCondition()
		buf := make([]byte, utf8.MaxRune+1)
		for r := rune(0); r <= utf8.MaxRune; r++ {
			buf[r] = byte(c.RuneWidth(r))
		}
		gotSHA := fmt.Sprintf("%x", sha256.Sum256(buf))
		if gotSHA != testcase.wantSHA {
			t.Errorf("TestRuneWidthChecksums = %s,\n\tsha256 = %s want %s",
				testcase.name, gotSHA, testcase.wantSHA)
		}
	}
}

func checkInterval(first, last rune) bool {
	return first >= 0 && first <= utf8.MaxRune &&
		last >= 0 && last <= utf8.MaxRune &&
		first <= last
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
	{'', 2},
	{'﬿', 1},
	{'\uFD37', 1},
	{'\U0001F1EB', 1},
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
	{"🇫🇮", 2},
	{"🏳️🌈", 4},
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
	s := "🇫🇮いうえお"
	expected := "🇫🇮いうえお"

	if out := Truncate(s, 10, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

func TestTruncate(t *testing.T) {
	s := "🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."
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
	s := "a🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "a🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおお..."

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
	s := "🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"
	expected := "🇫🇮いうえおあいうえおえおおおおおおおおおおおおおおおおおおおおおおおおおおおおお"

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
	s := "🇫🇮いうえおあい"
	expected := "🇫🇮いうえおあい"

	if out := Truncate(s, 80, "..."); out != expected {
		t.Errorf("Truncate(%q) = %q, want %q", s, out, expected)
	}
}

var truncatelefttests = []struct {
	s      string
	w      int
	prefix string
	out    string
}{
	{"source", 4, "", "ce"},
	{"source", 4, "...", "...ce"},
	{"あいうえお", 6, "", "えお"},
	{"あいうえお", 6, "...", "...えお"},
	{"あいうえお", 10, "", ""},
	{"あいうえお", 10, "...", "..."},
	{"あいうえお", 5, "", " えお"},
	{"Aあいうえお", 5, "", "うえお"},
	{"🇫🇮いうえお", 6, "", "えお"},
	{"🇫🇮いうえお", 6, "...", "...えお"},
	{"🇫🇮いうえお", 10, "", ""},
	{"🇫🇮いうえお", 10, "...", "..."},
	{"🇫🇮いうえお", 5, "", " えお"},
	{"A🇫🇮いうえお", 5, "", "うえお"},
}

func TestTruncateLeft(t *testing.T) {
	t.Parallel()

	for _, tt := range truncatelefttests {
		if out := TruncateLeft(tt.s, tt.w, tt.prefix); out != tt.out {
			t.Errorf("TruncateLeft(%q) = %q, want %q", tt.s, out, tt.out)
		}
	}
}

func TestFillLeft(t *testing.T) {
	s := "🇫🇮xいうえお"
	expected := "    🇫🇮xいうえお"

	if out := FillLeft(s, 15); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillLeftFit(t *testing.T) {
	s := "🇫🇮いうえお"
	expected := "🇫🇮いうえお"

	if out := FillLeft(s, 10); out != expected {
		t.Errorf("FillLeft(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRight(t *testing.T) {
	s := "🇫🇮xいうえお"
	expected := "🇫🇮xいうえお    "

	if out := FillRight(s, 15); out != expected {
		t.Errorf("FillRight(%q) = %q, want %q", s, out, expected)
	}
}

func TestFillRightFit(t *testing.T) {
	s := "🇫🇮いうえお"
	expected := "🇫🇮いうえお"

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
		{"👩\u200d", 2},
		{"👩\u200d🍳", 4},
		{"\u200d🍳", 2},
		{"👨\u200d👨", 4},
		{"👨\u200d👨\u200d👧", 6},
		{"🏳️\u200d🌈", 4},
		{"あ👩\u200d🍳い", 8},
		{"あ\u200d🍳い", 6},
		{"あ\u200dい", 4},
	}

	for _, tt := range tests {
		if got := c.StringWidth(tt.in); got != tt.want {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
