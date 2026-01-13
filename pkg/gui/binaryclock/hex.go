package binaryclock

import "gioui.org/layout"

type Hex struct {
	Val uint8
}

func (h Hex) Layout(gtx layout.Context) layout.Dimensions {
	_ = [4]Quad{
		{Val: (h.Val >> 6) & 0x03},
		{Val: (h.Val >> 4) & 0x03},
		{Val: (h.Val >> 2) & 0x03},
		{Val: (h.Val >> 0) & 0x03},
	}

	return layout.Dimensions{}
}
