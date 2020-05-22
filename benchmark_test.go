package runewidth

import (
	"testing"
	"unicode/utf8"
)

var benchSink int

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
