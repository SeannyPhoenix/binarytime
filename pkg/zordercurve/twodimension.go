package zordercurve

var TwoDimension twoDimension

type twoDimension struct{}

func (td twoDimension) ValidateCoord(c uint64) bool {
	return c&0x00000000ffffffff == c
}

func (td twoDimension) Dilate(v uint64) uint64 {
	v = (v | v<<16) & 0x0000ffff0000ffff
	v = (v | v<<8) & 0x00ff00ff00ff00ff
	v = (v | v<<4) & 0x0f0f0f0f0f0f0f0f
	v = (v | v<<2) & 0x3333333333333333
	v = (v | v<<1) & 0x5555555555555555

	return v
}

func (td twoDimension) Compress(v uint64) uint64 {
	v = v & 0x5555555555555555

	v = (v | v>>1) & 0x3333333333333333
	v = (v | v>>2) & 0x0f0f0f0f0f0f0f0f
	v = (v | v>>4) & 0x00ff00ff00ff00ff
	v = (v | v>>8) & 0x0000ffff0000ffff
	v = (v | v>>16) & 0x00000000ffffffff

	return v
}

func (td twoDimension) GetValue(x, y uint64) uint64 {
	return td.Dilate(x) | (td.Dilate(y) << 1)
}

func (td twoDimension) GetCoords(v uint64) (uint64, uint64) {
	return td.Compress(v), td.Compress(v >> 1)
}
