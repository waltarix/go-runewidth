package runewidth

import (
	"testing"
	"unicode/utf8"
)

var benchSink int

//
// RuneWidth
//

func benchRuneWidth(b *testing.B, start, stop rune, want int) int {
	n := 0
	got := -1
	c := NewCondition()
	for i := 0; i < b.N; i++ {
		got = n
		for r := start; r < stop; r++ {
			n += c.RuneWidth(r)
		}
		got = n - got
	}
	if want != 0 && got != want { // some extra checks
		b.Errorf("got %d, want %d\n", got, want)
	}
	return n
}
func BenchmarkRuneWidthAll(b *testing.B) {
	benchSink = benchRuneWidth(b, 0, utf8.MaxRune+1, 1430597)
}
func BenchmarkRuneWidth768(b *testing.B) {
	benchSink = benchRuneWidth(b, 0, 0x300, 715)
}

//
// String1Width - strings which consist of a single rune
//

func benchString1Width(b *testing.B, start, stop rune, want int) int {
	n := 0
	got := -1
	c := NewCondition()
	for i := 0; i < b.N; i++ {
		got = n
		for r := start; r < stop; r++ {
			s := string(r)
			n += c.StringWidth(s)
		}
		got = n - got
	}
	if want != 0 && got != want { // some extra checks
		b.Errorf("got %d, want %d\n", got, want)
	}
	return n
}
func BenchmarkString1WidthAll(b *testing.B) {
	benchSink = benchString1Width(b, 0, utf8.MaxRune+1, 1434693)
}
func BenchmarkString1Width768(b *testing.B) {
	benchSink = benchString1Width(b, 0, 0x300, 715)
}

//
// tables
//
func benchTable(b *testing.B) int {
	n := 0
	for i := 0; i < b.N; i++ {
		for r := rune(0); r <= utf8.MaxRune; r++ {
			if wcwidth9_width(r) >= -1 {
				n++
			}
		}
	}
	return n
}

func BenchmarkTableWcwidth9(b *testing.B) {
	benchSink = benchTable(b)
}
