package btclock

import (
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

type Clock struct {
	Granularity binarytime.Granularity

	// frames per second
	// zero means as soon as posible
	Framerate int
}

func (c Clock) Layout(gtx layout.Context, t binarytime.Date, theme *material.Theme) layout.Dimensions {
	label := material.H3(theme, t.HexGranular(c.Granularity))
	label.Font.Typeface = "monospace"
	label.Alignment = text.Middle
	label.Font.Weight = font.Bold
	label.Color = White
	label.MaxLines = 1

	f := min(max(c.Framerate, 0), 120)
	if f == 0 {
		gtx.Execute(op.InvalidateCmd{})
	} else {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / time.Duration(f))})
	}

	return label.Layout(gtx)
}

var (
	White = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	Black = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
)
