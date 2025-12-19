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
	bcapp "github.com/seannyphoenix/binarytime/internal/clock/app"
	"github.com/seannyphoenix/binarytime/internal/clock/window"
	"github.com/seannyphoenix/binarytime/pkg/binarytime"
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
	var state bcapp.State

	if len(os.Args) > 1 {
		if os.Args[1] == "-f" {
			state.ShowFrameRate = true
		}
	}

	window := window.New()
	theme := material.NewTheme()

	var t binarytime.Date

	startTime := time.Now()
	var frameCount int

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

			stack := op.Offset(image.Point{X: 0, Y: gtx.Dp(unit.Dp(24))}).Push(gtx.Ops)
			layoutTimeLabel(gtx, t, theme)
			stack.Pop()

			if !state.Sized {
				// Size the window to fit the current clock string plus some padding.
				// layout.Dimensions are in px; app.Size expects dp.
				const padDp = unit.Dp(24)
				// Re-measure the label to get its pixel dimensions.
				mgtx := gtx
				mgtx.Constraints = layout.Constraints{
					Min: image.Point{},
					Max: image.Pt(1<<30, 1<<30),
				}
				dims := layoutTimeLabel(mgtx, t, theme)

				// Convert px -> dp using the current metric.
				pxPerDp := gtx.Metric.PxPerDp
				wantW := unit.Dp(float32(dims.Size.X)/pxPerDp) + padDp*2
				wantH := unit.Dp(float32(dims.Size.Y)/pxPerDp) + padDp*2

				window.Option(app.Size(wantW, wantH))
				state.Sized = true
			}

			if state.ShowFrameRate {
				layoutCountLabel(gtx, theme, frameCount, avgFpS)
			}

			gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(100 * time.Millisecond)})

			e.Frame(gtx.Ops)
		}
	}
}

func layoutTimeLabel(gtx layout.Context, t binarytime.Date, theme *material.Theme) layout.Dimensions {
	label := material.H3(theme, t.String())
	label.Font.Typeface = "monospace"
	label.Font.Weight = font.Bold
	label.Alignment = text.Middle
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	label.MaxLines = 1
	return label.Layout(gtx)
}

func layoutCountLabel(gtx layout.Context, theme *material.Theme, count int, avgFpS float64) layout.Dimensions {
	label := material.Caption(theme, fmt.Sprintf("Frame count: %d, Avg FPS: %.2f", count, avgFpS))
	label.Font.Typeface = "monospace"
	label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	return label.Layout(gtx)
}
