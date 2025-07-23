package byteglyph

func Glyph(v byte) string {
	return newHorizontalGlyph(v).string()
}

func Glyphs(vs []byte, dot int) string {
	hgs := make([]horizontalGlyph, len(vs))
	for i, v := range vs {
		hgs[i] = newHorizontalGlyph(v)
	}

	return assembleHorizontalGlyphs(hgs, dot)
}
