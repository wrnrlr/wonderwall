package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
)

var red = color.NRGBA{R: 255, A: 255}

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				// Draw a red rectangle.
				stack := op.Save(&ops)
				clip.Rect{Max: image.Pt(100, 50)}.Add(&ops)
				paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(&ops)
				paint.PaintOp{}.Add(&ops)
				stack.Load()

				// Draw a green rectangle.
				stack = op.Save(&ops)
				clip.Rect{Max: image.Pt(50, 100)}.Add(&ops)
				paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(&ops)
				paint.PaintOp{}.Add(&ops)
				stack.Load()

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
