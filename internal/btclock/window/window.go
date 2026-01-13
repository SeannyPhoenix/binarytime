package window

import "gioui.org/app"

type Window struct {
	app.Window
}

func New(options ...app.Option) *app.Window {
	w := &app.Window{}
	options = append(options, app.Title("binarytime"))
	w.Option(options...)
	return w
}
