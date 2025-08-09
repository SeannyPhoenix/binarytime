package binarytime

import (
	"encoding/binary"
	"fmt"

	"github.com/seannyphoenix/binarytime/pkg/byteglyph"
)

func (bt Date) String() string {
	return coarse(bt)
}

func (bt Date) StringFine() string {
	return fine(bt)
}

// Glyphs returns a string representation of the BinaryTime using byteglyphs.
// It uses all 128 bits: the first 8 bytes for the date and the next 8 bytes
// for the time.
func (bt Date) Glyphs() string {
	return byteglyph.Glyphs(bt.Bytes(), 8)
}

// TimeGlyphs returns a string representation of the time portion of the BinaryTime
// using byteglyphs. Only the 8th and 9th bytes are used, which represent the time
// down to the seconds level.
func (bt Date) TimeGlyphs() string {
	return byteglyph.Glyphs(bt.Bytes()[8:10], 0)
}

// DateGlyphs returns a string representation of the date portion of the BinaryTime
// using byteglyphs. It uses the 6th and 7th bytes, which represent the date up to
// the centuries level.
func (bt Date) DateGlyphs() string {
	return byteglyph.Glyphs(bt.Bytes()[6:8], 2)
}

// DateTimeGlyphs returns a string representation of the date and time portion
// of the BinaryTime using byteglyphs. It uses the 6th through 9th bytes,
// which represent the date up to the centuries level and time down to the seconds level.
func (bt Date) DateTimeGlyphs() string {
	return byteglyph.Glyphs(bt.Bytes()[6:10], 2)
}

func coarse(bt Date) string {
	bytes := bt.Bytes()

	days := binary.BigEndian.Uint16(bytes[6:8])
	subDays := binary.BigEndian.Uint16(bytes[8:10])

	return fmt.Sprintf("%04X:%04X", days, subDays)
}

func fine(bt Date) string {
	bytes := bt.Bytes()
	days := binary.BigEndian.Uint64(bytes[:8])
	subDays := binary.BigEndian.Uint64(bytes[8:])

	return fmt.Sprintf("%016X:%016X", days, subDays)
}
