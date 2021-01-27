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
	"github.com/wrnrlr/wonderwall/wonder/f32color"
	"image"
	"image/color"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		quarter := uint8(0x40)
		primary := f32color.RGBAToNRGBA(color.RGBA{G: 0xff, A: 0xff})
		//white := f32color.RGBAToNRGBA(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: quarter})
		black := f32color.RGBAToNRGBA(color.RGBA{A: quarter})
		transparent := f32color.RGBAToNRGBA(color.RGBA{A: 0xff})
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				dr := image.Rectangle{Max: gtx.Constraints.Min}
				stack := op.Save(gtx.Ops)
				topRight := f32.Point{X: float32(dr.Max.X), Y: float32(dr.Min.Y)}
				bottomLeft := f32.Point{X: float32(dr.Min.X), Y: float32(dr.Max.Y)}
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
					Color1: color.NRGBA{A: 0xff},
					Color2: color.NRGBA{R: 0x00, G: 0x00, B: 0x0, A: 0xff},
				}.Add(gtx.Ops)
				clip.Rect(dr).Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				paint.LinearGradientOp{
					Stop1:  topLeft,
					Stop2:  bottomLeft,
					Color1: transparent,
					Color2: black,
				}.Add(gtx.Ops)
				//clip.Rect(dr).Add(gtx.Ops)
				//paint.PaintOp{}.Add(gtx.Ops)
				stack.Load()
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
