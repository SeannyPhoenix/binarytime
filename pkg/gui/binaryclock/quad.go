package binaryclock

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Quad struct {
	Val uint8
}

// increment increases the quad's value by one,
// base 4. It returns true if the value carries
func (q *Quad) Increment() bool {
	q.Val = (q.Val + 1) % 4
	return q.Val == 0
}

func (q Quad) Layout(gtx layout.Context) layout.Dimensions {
	// Create a 2x2 grid of squares. Compute a fixed square size so the
	// top and bottom rows touch precisely (no extra gap introduced by
	// flex allocations). The quad will be 2*s by 2*s.

	avail := gtx.Constraints.Max

	// Cap the maximum square size to avoid huge squares on large displays.
	// const maxDp = 48
	// maxPx := gtx.Dp(unit.Dp(maxDp))

	// Candidate square size is limited by half the available width and
	// the available height (since we need two rows). Choose the minimum
	// of those and the DP cap.
	s := max(min(avail.Y, avail.X/2), 0)

	total := image.Pt(2*s, 2*s)
	// Ensure the quad doesn't request more than the parent can give.
	if total.X > avail.X {
		total.X = avail.X
	}
	if total.Y > avail.Y {
		total.Y = avail.Y
	}

	// Constrain ourselves to the computed total size.
	gtx.Constraints.Max = total

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Top row: two rigid squares of size s
			rowSize := image.Pt(2*s, s)
			gtx.Constraints.Max = rowSize
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return q.drawSquare(gtx, q.Val >= 1, image.Pt(s, s))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return q.drawSquare(gtx, q.Val >= 2, image.Pt(s, s))
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Bottom row: two rigid squares of size s
			rowSize := image.Pt(2*s, s)
			gtx.Constraints.Max = rowSize
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return q.drawSquare(gtx, q.Val >= 3, image.Pt(s, s))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return q.drawSquare(gtx, false, image.Pt(s, s))
				}),
			)
		}),
	)
}

// drawSquare draws a square, filled or empty based on the filled parameter
// drawSquare draws a square of the explicit size provided so the caller can
// control exact positioning. This avoids flex-introduced gaps between rows.
func (q *Quad) drawSquare(gtx layout.Context, filled bool, size image.Point) layout.Dimensions {
	// Constrain to the explicit size
	gtx.Constraints.Max = size

	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color.NRGBA{R: 50, G: 50, B: 50, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	inset := op.Offset(image.Point{X: gtx.Dp(unit.Dp(2)), Y: gtx.Dp(unit.Dp(2))}).Push(gtx.Ops)

	// Choose color based on fill state
	fillColor := color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Light gray for empty
	if filled {
		fillColor = color.NRGBA{R: 100, G: 40, B: 40, A: 255}
	}
 
	paint.ColorOp{Color: fillColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	inset.Pop()

	return layout.Dimensions{Size: size}
}
