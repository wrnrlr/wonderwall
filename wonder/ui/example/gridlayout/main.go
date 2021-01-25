package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/ui"
	"image/color"
)

var red = color.NRGBA{R: 255, A: 255}

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		th := material.NewTheme(gofont.Collection())
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				ui.Grid{Columns: 5, Rows: 5}.Layout(gtx, func(i, j int, gtx ui.C) {
					material.Label(th, unit.Sp(16), fmt.Sprintf("%d", i*5+j)).Layout(gtx)
				})
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
