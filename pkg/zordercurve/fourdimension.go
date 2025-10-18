package zordercurve

type FourDimension struct {
}

func (z FourDimension) ValidateCoord(c uint64) bool {
	return c&0x00000000000ffff == c
}
