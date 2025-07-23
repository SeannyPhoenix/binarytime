package byteglyph

import (
	"fmt"
	"strings"
)

var (
	emptyH = []byte(`    `)
	highH  = []byte(`|\/|`)
	barH   = []byte(`----`)
	lowH   = []byte(`|/\|`)
	dotH   = []byte("*   ")
)

type horizontalGlyph struct {
	high []byte
	bar  []byte
	low  []byte
}

func newHorizontalGlyph(b byte) horizontalGlyph {
	hg := horizontalGlyph{
		high: make([]byte, 4),
		bar:  make([]byte, 4),
		low:  make([]byte, 4),
	}
	copy(hg.high, emptyH)
	copy(hg.bar, barH)
	copy(hg.low, emptyH)

	l := byte(0b10000000)
	for i := range 4 {
		if b&l != 0 {
			hg.high[i] = highH[i]
		}
		l >>= 1
	}

	for i := range 4 {
		if b&l != 0 {
			hg.low[i] = lowH[i]
		}
		l >>= 1
	}

	return hg
}

func (hg horizontalGlyph) string() string {
	return fmt.Sprintf("%s\n%s\n%s\n", hg.high, hg.bar, hg.low)
}

func assembleHorizontalGlyphs(hgs []horizontalGlyph, dot int) string {
	var sb strings.Builder
	for i, hg := range hgs {
		if i == dot {
			sb.Write(dotH)
			sb.WriteString("\n\n")
		}
		sb.WriteString(hg.string())
		sb.WriteRune('\n')
	}
	if dot == len(hgs) {
		sb.Write(dotH)
		sb.WriteString("\n\n")
	}
	return sb.String()
}
