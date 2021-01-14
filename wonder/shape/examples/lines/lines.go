package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image/color"
)

var red = color.NRGBA{R: 255}

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				shape.NewPolyline([]f32.Point{{0, 0}, {100, 100}}, red, 10).Draw(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
