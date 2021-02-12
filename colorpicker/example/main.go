package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/colorpicker"
)

var th *material.Theme

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		th = material.NewTheme(gofont.Collection())
		colorPicker := colorpicker.NewColorPicker(th)
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				colorPicker.Event()
				colorPicker.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
