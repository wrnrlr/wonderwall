package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/wrnrlr/wonderwall/wonder/nexttool"
)

var th *material.Theme

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		th = material.NewTheme(gofont.Collection())
		toolbar := nexttool.NewToolMenu(th,
			&nexttool.Arrange{},
			&nexttool.Selection{},
			&nexttool.Brush{},
			&nexttool.Text{},
			&nexttool.Shape{},
			&nexttool.Image{},
			&nexttool.Pen{},
			&nexttool.Zoom{})
		var ops op.Ops
		for {
			ev := <-w.Events()
			switch e := ev.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				bg := layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
				fg := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return toolbar.Layout(gtx)
				})
				layout.Stack{Alignment: layout.E}.Layout(gtx, bg, fg)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
