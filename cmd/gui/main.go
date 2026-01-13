package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/seannyphoenix/binarytime/internal/btclock"
	"github.com/seannyphoenix/binarytime/internal/btclock/btapp"
	"github.com/seannyphoenix/binarytime/internal/btclock/window"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
	"github.com/seannyphoenix/binarytime/pkg/gui/binaryclock"
)

func main() {
	go func() {
		err := run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run() error {
	var state btapp.State

	if len(os.Args) > 1 {
		if os.Args[1] == "-f" {
			state.ShowFrameRate = true
		}
	}

	window := window.New(app.Decorated(false))
	theme := material.NewTheme()

	var t binarytime.Date

	startTime := time.Now()
	var frameCount int

	c := btclock.Clock{
		Framerate: 25,
	}

	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			frameCount++
			avgFpS := float64(frameCount) / time.Since(startTime).Seconds()

			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 0xff})

			t = binarytime.DateFromTime(e.Now)

			// if !state.Sized {
			// 	// Size the window to fit the current clock string plus some padding.
			// 	// layout.Dimensions are in px; app.Size expects dp.
			// 	const padDp = unit.Dp(24)
			// 	// Re-measure the label to get its pixel dimensions.
			// 	mgtx := gtx
			// 	mgtx.Constraints = layout.Constraints{
			// 		Min: image.Point{},
			// 		Max: image.Pt(1<<30, 1<<30),
			// 	}
			// 	dims := c.Layout(mgtx, t, theme)

			// 	// Convert px -> dp using the current metric.
			// 	pxPerDp := gtx.Metric.PxPerDp
			// 	wantW := unit.Dp(float32(dims.Size.X)/pxPerDp) + padDp*2
			// 	wantH := unit.Dp(float32(dims.Size.Y)/pxPerDp) + padDp*2

			// 	window.Option(app.Size(wantW, wantH))
			// 	state.Sized = true
			// }

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					stack := op.Offset(image.Point{X: 0, Y: gtx.Dp(unit.Dp(24))}).Push(gtx.Ops)
					c.Layout(gtx, t, theme)
					stack.Pop()
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Flexed(1,
					binaryclock.Clock{Time: t}.Layout,
				),
			)

			if state.ShowFrameRate {
				layoutCountLabel(gtx, theme, frameCount, avgFpS)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func layoutCountLabel(gtx layout.Context, theme *material.Theme, count int, avgFpS float64) layout.Dimensions {
	flex := layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}

	return flex.Layout(gtx, layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			label := material.Caption(theme, fmt.Sprintf("Frame count: %d", count))
			label.Font.Typeface = "monospace"
			label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
			return label.Layout(gtx)
		},
	),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				label := material.Caption(theme, fmt.Sprintf("Avg FPS: %.2f", avgFpS))
				label.Font.Typeface = "monospace"
				label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
				return label.Layout(gtx)
			},
		),
	)
}
