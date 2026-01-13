package binarytime

import (
	"github.com/seannyphoenix/binarytime/pkg/byteglyph"
)

func (d Date) MarshalText() ([]byte, error) {
	return d.value.MarshalText()
}

func (d Date) String() string {
	return d.Hex()
}

func (d Date) StringFine() string {
	return d.HexFine()
}

func (d Date) Hex() string {
	return coarse(d)
}

type Granularity struct {
	Upper int
	Lower int
}

func (g Granularity) ul() (int, int) {
	if g.Upper <= 0 {
		g.Upper = 2
	}
	g.Upper = min(g.Upper, 8)
	if g.Lower <= 0 {
		g.Lower = 2
	}
	g.Lower = min(g.Lower, 8)

	return 9 - g.Upper, 9 + g.Lower
}

func (d Date) HexGranular(g Granularity) string {
	s, _ := d.value.StringWithPrecision(g.ul())
	return s
}

func (d Date) HexFine() string {
	return fine(d)
}

func (d Date) Base64() string {
	return d.value.Base64()
}

// Glyphs returns a string representation of the BinaryTime using byteglyphs.
// It uses all 128 bits: the first 8 bytes for the date and the next 8 bytes
// for the time.
func (d Date) Glyphs() string {
	return byteglyph.Glyphs(d.Bytes(), 8)
}

// TimeGlyphs returns a string representation of the time portion of the BinaryTime
// using byteglyphs. Only the 8th and 9th bytes are used, which represent the time
// down to the seconds level.
func (d Date) TimeGlyphs() string {
	return byteglyph.Glyphs(d.Bytes()[8:10], 0)
}

// DateGlyphs returns a string representation of the date portion of the BinaryTime
// using byteglyphs. It uses the 6th and 7th bytes, which represent the date up to
// the centuries level.
func (d Date) DateGlyphs() string {
	return byteglyph.Glyphs(d.Bytes()[6:8], 2)
}

// DateTimeGlyphs returns a string representation of the date and time portion
// of the BinaryTime using byteglyphs. It uses the 6th through 9th bytes,
// which represent the date up to the centuries level and time down to the seconds level.
func (d Date) DateTimeGlyphs() string {
	return byteglyph.Glyphs(d.Bytes()[6:10], 2)
}

func coarse(d Date) string {
	s, _ := d.value.StringWithPrecision(7, 11)
	return s
}

func fine(d Date) string {
	return d.value.String()
}
