package clock

import "gioui.org/layout"

type Granularity struct {
	Upper int
	Lower int
}

type Clock struct {
	Granularity Granularity
}

func (Clock) Layout() layout.Dimensions {
	return layout.Dimensions{}
}
