package binarytime

import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"

	"github.com/seannyphoenix/binarytime/pkg/byteglyph"
	"github.com/seannyphoenix/binarytime/pkg/fixed128"
)

var (
	BinaryTimeRegexp           = regexp.MustCompile(`^@[0-9a-f]{16}\.[0-9a-f]{16}$`)
	ErrInvalidBinaryTimeFormat = errors.New("invalid binary time format")
)

var (
	_ encoding.BinaryMarshaler   = (*Date)(nil)
	_ encoding.BinaryUnmarshaler = (*Date)(nil)

	_ encoding.TextMarshaler   = (*Date)(nil)
	_ encoding.TextUnmarshaler = (*Date)(nil)
)

func (d Date) MarshalText() ([]byte, error) {
	hi := d.value.Bytes()[0:8]
	lo := d.value.Bytes()[8:16]

	return fmt.Appendf(nil, "@%016x.%016x", hi, lo), nil
}

func (d *Date) UnmarshalText(text []byte) error {
	if !BinaryTimeRegexp.Match(text) {
		return fmt.Errorf("%w: %s", ErrInvalidBinaryTimeFormat, text)
	}

	hi := binary.BigEndian.Uint64(text[1:17])
	lo := binary.BigEndian.Uint64(text[18:34])

	d.value = fixed128.FromParts(hi, lo, false)
	return nil
}

func (d Date) MarshalBinary() ([]byte, error) {
	return d.value.Bytes(), nil
}

func (d *Date) UnmarshalBinary(data []byte) error {
	value, err := fixed128.FromBytes(data)
	if err != nil {
		return err
	}
	d.value = value
	return nil
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
