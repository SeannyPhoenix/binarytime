package zordercurve

var ThreeDimension threeDimension

type threeDimension struct{}

func (td threeDimension) ValidateCoord(c uint64) bool {
	return c&0x00000000001fffff == c
}

func (td threeDimension) Dilate(v uint64) uint64 {
	return v
}

func (td threeDimension) Compress(v uint64) uint64 {
	return v
}

func (td threeDimension) GetValue(x, y, z uint64) uint64 {
	return td.Dilate(x) | (td.Dilate(y) << 1) | (td.Dilate(z) << 2)
}

func (td threeDimension) GetCoords(v uint64) (uint64, uint64, uint64) {
	return td.Compress(v), td.Compress(v >> 1), td.Compress(v >> 2)
}
