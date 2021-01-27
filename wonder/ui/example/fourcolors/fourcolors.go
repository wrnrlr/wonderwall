package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/wrnrlr/wonderwall/wonder/f32color"
	"image"
	"image/color"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		//quarter := uint8(0x40)
		colorPicker := &ColorPicker{}
		//white := f32color.RGBAToNRGBA(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: quarter})
		//black := f32color.RGBAToNRGBA(color.RGBA{A: quarter})
		//transparent := f32color.RGBAToNRGBA(color.RGBA{A: 0xff})
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				colorPicker.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}

type ColorPicker struct {
	input   widget.Editor
	rainbow widget.Float
	alfa    widget.Float
}

var primary = f32color.RGBAToNRGBA(color.RGBA{G: 0xff, A: 0xff})

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(cp.layoutGradiants))
	})
}

func (cp *ColorPicker) layoutGradiants(gtx layout.Context) layout.Dimensions {
	dr := image.Rectangle{Max: image.Point{X: 200, Y: 150}}
	stack := op.Save(gtx.Ops)
	topRight := f32.Point{X: float32(dr.Max.X), Y: float32(dr.Min.Y)}
	//bottomLeft := f32.Point{X: float32(dr.Min.X), Y: float32(dr.Max.Y)}
	topLeft := f32.Point{X: float32(dr.Min.X), Y: float32(dr.Min.Y)}
	bottomRight := f32.Point{X: float32(dr.Max.X), Y: float32(dr.Max.Y)}
	paint.LinearGradientOp{
		Stop1:  topRight,
		Stop2:  topLeft,
		Color1: primary,
		Color2: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}.Add(gtx.Ops)
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	paint.LinearGradientOp{
		Stop1:  topRight,
		Stop2:  bottomRight,
		Color1: color.NRGBA{},
		Color2: color.NRGBA{R: 0x00, G: 0x00, B: 0x0, A: 0xff},
	}.Add(gtx.Ops)
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	stack.Load()
	return layout.Dimensions{Size: dr.Max}
}

func (cp *ColorPicker) Events(gtx layout.Context) {

}
