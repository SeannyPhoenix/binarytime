package fixed128

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const signIndex = 0

var (
	ErrInvalidLength  = errors.New("invalid length")
	ErrInvalidSignBit = errors.New("invalid sign bit")

	ErrInvalidFormat = errors.New("invalid format")
)

// String implements fmt.Stringer
func (f128 Fixed128) String() string {
	if f128.IsZero() {
		return "00.00"
	}

	b := f128.bytes()

	high := 1
	for high < 8 && b[high] == 0 {
		high++
	}

	low := 17
	for low > 10 && b[low-1] == 0 {
		low--
	}

	var out string
	if f128.IsNeg() {
		out = "-"
	}

	out += fmt.Sprintf("%0*X.%0*X", (9-high)*2, b[high:9], (low-9)*2, b[9:low])
	return out
}

func Parse(s string) (Fixed128, error) {
	var f128 Fixed128
	if err := f128.UnmarshalText([]byte(s)); err != nil {
		return Fixed128{}, err
	}
	return f128, nil
}

// MarshalBinary implements encoding.BinaryMarshaler
func (f128 Fixed128) MarshalBinary() ([]byte, error) {
	return f128.bytes(), nil
}

func (f128 Fixed128) bytes() []byte {
	var b [17]byte
	f128.value.FillBytes(b[1:17])

	if f128.IsNeg() {
		b[signIndex] = 1
	}

	return b[:]
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (f128 *Fixed128) UnmarshalBinary(bin []byte) error {
	if len(bin) != 17 {
		return fmt.Errorf("%w: got %d", ErrInvalidLength, len(bin))
	}
	if bin[signIndex] != 0 && bin[signIndex] != 1 {
		return fmt.Errorf("%w: got %02X", ErrInvalidSignBit, bin[signIndex])
	}

	f128.value.SetBytes(bin[1:])

	if bin[signIndex] == 1 {
		f128.value.Neg(&f128.value)
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (f128 Fixed128) MarshalText() ([]byte, error) {
	return []byte(f128.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// Accepts the same format produced by String(): e.g. "-03FC65D0.1003F0".
func (f128 *Fixed128) UnmarshalText(text []byte) error {
	s := string(text)

	if len(s) == 0 {
		return fmt.Errorf("%w: empty string", ErrInvalidFormat)
	}
	if len(s) > 34 {
		return fmt.Errorf("%w: too long", ErrInvalidFormat)
	}

	neg := s[0] == '-'
	if neg {
		s = s[1:]
	}

	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return fmt.Errorf("%w: expected 'HI.LO'", ErrInvalidFormat)
	}

	// Parse hex big-endian whole+frac into 16 bytes
	hiStr, loStr := parts[0], parts[1]

	hiBytes, err := hex.DecodeString(hiStr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	loBytes, err := hex.DecodeString(loStr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}

	if len(hiBytes) > 8 || len(loBytes) > 8 {
		return fmt.Errorf("%w: too wide", ErrInvalidFormat)
	}

	var buf [17]byte
	copy(buf[1+(8-len(hiBytes)):9], hiBytes)
	copy(buf[9:9+len(loBytes)], loBytes)
	if neg {
		buf[signIndex] = 1
	}
	return f128.UnmarshalBinary(buf[:])
}
