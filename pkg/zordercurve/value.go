package zordercurve

func GetValueFromXY(x, y uint32) uint64 {
	return dilate(x) | (dilate(y) << 1)
}

func GetXYFromValue(v uint64) (uint32, uint32) {
	return compress(v), compress(v >> 1)
}

func dilate(v uint32) uint64 {
	d := uint64(v)

	d = (d | d<<16) & 0x0000ffff0000ffff
	d = (d | d<<8) & 0x00ff00ff00ff00ff
	d = (d | d<<4) & 0x0f0f0f0f0f0f0f0f
	d = (d | d<<2) & 0x3333333333333333
	d = (d | d<<1) & 0x5555555555555555

	return d
}

func compress(v uint64) uint32 {
	v = v & 0x5555555555555555

	v = (v | v>>1) & 0x3333333333333333
	v = (v | v>>2) & 0x0f0f0f0f0f0f0f0f
	v = (v | v>>4) & 0x00ff00ff00ff00ff
	v = (v | v>>8) & 0x0000ffff0000ffff
	v = (v | v>>16) & 0x00000000ffffffff

	return uint32(v)
}
