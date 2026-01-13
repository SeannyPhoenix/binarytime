package binaryclock

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

type Clock struct {
	Time binarytime.Date
}

func (c Clock) Layout(gtx layout.Context) layout.Dimensions {
	quads := [16]uint8{
		c.Time.Bytes()[6] >> 6 & 0x03,
		c.Time.Bytes()[6] >> 4 & 0x03,
		c.Time.Bytes()[6] >> 2 & 0x03,
		c.Time.Bytes()[6] & 0x03,
		c.Time.Bytes()[7] >> 6 & 0x03,
		c.Time.Bytes()[7] >> 4 & 0x03,
		c.Time.Bytes()[7] >> 2 & 0x03,
		c.Time.Bytes()[7] & 0x03,
		c.Time.Bytes()[8] >> 6 & 0x03,
		c.Time.Bytes()[8] >> 4 & 0x03,
		c.Time.Bytes()[8] >> 2 & 0x03,
		c.Time.Bytes()[8] & 0x03,
		c.Time.Bytes()[9] >> 6 & 0x03,
		c.Time.Bytes()[9] >> 4 & 0x03,
		c.Time.Bytes()[9] >> 2 & 0x03,
		c.Time.Bytes()[9] & 0x03,
	}

	// spacing between quads, scaled by display density
	spacer := func() layout.FlexChild {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Width: unit.Dp(6)}.Layout(gtx)
		})
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[0]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[1]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[2]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[3]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[4]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[5]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[6]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[7]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[8]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[9]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[10]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[11]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[12]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[13]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[14]}; return cl.Layout(gtx) }), spacer(),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { cl := Quad{Val: quads[15]}; return cl.Layout(gtx) }),
	)
}
