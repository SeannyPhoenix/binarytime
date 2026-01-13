package main

import (
	"image"
	"image/color"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/seannyphoenix/binarytime/pkg/timer"
)

func main() {
	go run()
	app.Main()
}

func run() {
	w := &app.Window{}
	w.Option(
		app.Title("Egg timer"),
		app.Size(unit.Dp(400), unit.Dp(600)),
	)

	var ops op.Ops
	th := material.NewTheme()

	var startButton widget.Clickable

	var timer timer.Timer
	timer.Set(5 * time.Second)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			var btn material.ButtonStyle

			if startButton.Clicked(gtx) {
				if timer.Finished() {
					timer.Reset()
				} else {
					timer.Toggle(e.Now)
				}
			}

			if timer.Running() {
				btn = material.Button(th, &startButton, "Stop")
			} else {
				if timer.Finished() {
					btn = material.Button(th, &startButton, "Reset")
				} else {
					btn = material.Button(th, &startButton, "Start")
				}
			}

			pb := material.ProgressBar(th, 0)

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					circle := clip.Ellipse{
						Min: image.Pt(gtx.Constraints.Max.X/2-120, 0),
						Max: image.Pt(gtx.Constraints.Max.X/2+120, 240),
					}.Op(gtx.Ops)
					paint.FillShape(gtx.Ops, color.NRGBA{
						R: 0xa0,
						G: uint8(0xa0 * (1 - timer.Progress())),
						B: uint8(0xa0 * (1 - timer.Progress())),
						A: 0xff,
					}, circle)
					return layout.Dimensions{Size: image.Point{Y: 400}}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if timer.Running() {
						timer.Tick(e.Now)
						gtx.Execute(op.InvalidateCmd{At: e.Now.Add(time.Second / 25)})
					}
					pb.Progress = timer.Progress()
					return pb.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(24)).Layout(gtx, btn.Layout)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
