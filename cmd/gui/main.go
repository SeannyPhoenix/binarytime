package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func run(window *app.Window) error {
	var sized, showFrameRate bool

	if len(os.Args) > 1 {
		if os.Args[1] == "-f" {
			showFrameRate = true
		}
	}

	window.Option(
		app.Title("Date"),
	)
	theme := material.NewTheme()

	var t binarytime.Date
	var c int

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			c++
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 0xff})

			t = binarytime.DateFromTime(e.Now)

			stack := op.Offset(image.Point{X: 0, Y: gtx.Dp(unit.Dp(24))}).Push(gtx.Ops)
			layoutTimeLabel(gtx, t, theme)
			stack.Pop()

			if !sized {
				// Size the window to fit the current clock string plus some padding.
				// layout.Dimensions are in px; app.Size expects dp.
				const padDp = unit.Dp(24)
				// Re-measure the label to get its pixel dimensions.
				mgtx := gtx
				mgtx.Constraints = layout.Constraints{
					Min: image.Point{},
					Max: image.Pt(1<<30, 1<<30),
				}
				dims := timeLabelStyle(theme, t).Layout(mgtx)

				// Convert px -> dp using the current metric.
				pxPerDp := gtx.Metric.PxPerDp
				wantW := unit.Dp(float32(dims.Size.X)/pxPerDp) + padDp*2
				wantH := unit.Dp(float32(dims.Size.Y)/pxPerDp) + padDp*2

				window.Option(app.Size(wantW, wantH))
				sized = true
			}

			if showFrameRate {
				countLabel(gtx, c, theme)
			}

			gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(100 * time.Millisecond)})

			e.Frame(gtx.Ops)
		}
	}
}

func timeLabelStyle(theme *material.Theme, t binarytime.Date) material.LabelStyle {
	label := material.H3(theme, t.String())
	label.Font.Typeface = "monospace"
	label.Font.Weight = font.Bold
	label.Alignment = text.Middle
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	label.MaxLines = 1
	return label
}

func layoutTimeLabel(gtx layout.Context, t binarytime.Date, theme *material.Theme) layout.Dimensions {
	return timeLabelStyle(theme, t).Layout(gtx)
}

func countLabel(gtx layout.Context, count int, theme *material.Theme) layout.Dimensions {
	label := material.Body1(theme, fmt.Sprintf("Frame count: %d", count))
	label.Font.Typeface = "monospace"
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	return label.Layout(gtx)
}
