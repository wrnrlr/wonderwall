package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/colorpicker"
	"image/color"
)

var th *material.Theme

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		th = material.NewTheme(gofont.Collection())
		colorField := colorpicker.NewColorSelection(th,
			colorpicker.NewPicker(th),
			colorpicker.NewAlphaSlider(),
			colorpicker.NewToggle(&widget.Clickable{},
				colorpicker.NewHexEditor(th),
				colorpicker.NewRgbEditor(th),
				colorpicker.NewHsvEditor(th)))
		colorField.SetColor(color.NRGBA{G: 255, A: 255})
		colorEditors := colorpicker.NewMux(
			colorpicker.NewPicker(th),
			colorpicker.NewHexEditor(th),
			colorpicker.NewRgbEditor(th),
			colorpicker.NewHsvEditor(th),
			colorpicker.NewAlphaSlider())
		colorEditors.SetColor(color.NRGBA{B: 255, A: 255})
		btn := &widget.Clickable{}
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				colorField.Changed()
				colorEditors.Changed()
				gtx := layout.NewContext(&ops, e)
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.H2(th, "Color Inputs").Layout),
					layout.Rigid(material.H3(th, "ColorSelection").Layout),
					layout.Rigid(material.Body1(th, "Click on color fields to show and hide a colorpicker.").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(4)).Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
									layout.Rigid(material.Label(th, unit.Sp(16), " Color: ").Layout),
									layout.Rigid(colorField.Layout),
									layout.Rigid(material.Button(th, btn, "Button").Layout))
							})
					}),
					layout.Rigid(material.H3(th, "Picker").Layout),
					layout.Rigid(material.Body1(th, "The color picker component can be used on its own.").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return colorEditors.Layout(gtx)
					}))
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
