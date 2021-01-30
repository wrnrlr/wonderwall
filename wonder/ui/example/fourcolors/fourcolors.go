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
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(cp.layoutGradiants),
			layout.Rigid(cp.layoutRainbow))
	})
}

func (cp *ColorPicker) layoutGradiants(gtx layout.Context) layout.Dimensions {
	w := gtx.Px(unit.Dp(200))
	h := gtx.Px(unit.Dp(120))
	dr := image.Rectangle{Max: image.Point{X: w, Y: h}}
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

	//p := f32.Point{X: 100, Y: 75}
	//drawCircle(p, gtx)

	stack.Load()
	return layout.Dimensions{Size: dr.Max}
}

var (
	red     = color.NRGBA{R: 255, A: 255}
	yellow  = color.NRGBA{R: 255, G: 255, A: 255}
	green   = color.NRGBA{G: 255, A: 255}
	cyan    = color.NRGBA{G: 255, B: 255, A: 255}
	skyblue = color.NRGBA{G: 127, B: 255, A: 255}
	blue    = color.NRGBA{B: 255, A: 255}
	purple  = color.NRGBA{R: 127, G: 0, B: 255, A: 255}
	magenta = color.NRGBA{R: 255, B: 255, A: 255}
)

var colors = []color.NRGBA{red, yellow, green, cyan, skyblue, blue, purple, magenta, red}

func (cp *ColorPicker) layoutRainbow(gtx layout.Context) layout.Dimensions {
	w := gtx.Px(unit.Dp(200))
	h := gtx.Px(unit.Dp(20))

	w2 := w / (len(colors) - 1)
	offsetX := 0
	color1 := colors[0]
	for _, color2 := range colors[1:] {
		stack := op.Save(gtx.Ops)
		paint.LinearGradientOp{
			Stop1:  f32.Point{float32(offsetX), 0},
			Stop2:  f32.Point{float32(offsetX + w2), 0},
			Color1: color1,
			Color2: color2,
		}.Add(gtx.Ops)
		dr := image.Rectangle{Min: image.Point{X: offsetX, Y: 0}, Max: image.Point{X: offsetX + w2, Y: h}}
		clip.Rect(dr).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()
		offsetX += w2
		color1 = color2
	}
	return layout.Dimensions{Size: image.Point{X: w, Y: h}}
}

const c = 0.55228475 // 4*(sqrt(2)-1)/3

func drawCircle(p f32.Point, gtx layout.Context) {
	width := float32(2)
	r := float32(5)
	w, h := r*2, r*2
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Point{X: p.X, Y: p.Y})
	path.Move(f32.Point{X: w / 4 * 3, Y: r / 2})
	path.Cube(f32.Point{X: 0, Y: r * c}, f32.Point{X: -r + r*c, Y: r}, f32.Point{X: -r, Y: r})    // SE
	path.Cube(f32.Point{X: -r * c, Y: 0}, f32.Point{X: -r, Y: -r + r*c}, f32.Point{X: -r, Y: -r}) // SW
	path.Cube(f32.Point{X: 0, Y: -r * c}, f32.Point{X: r - r*c, Y: -r}, f32.Point{X: r, Y: -r})   // NW
	path.Cube(f32.Point{X: r * c, Y: 0}, f32.Point{X: r, Y: r - r*c}, f32.Point{X: r, Y: r})      // NE
	path.Move(f32.Point{X: -w, Y: -r})                                                            // Return to origin
	scale := (r - width*2) / r
	path.Move(f32.Point{X: w * (1 - scale) * .5, Y: h * (1 - scale) * .5})
	w *= scale
	h *= scale
	r *= scale
	path.Move(f32.Point{X: 0, Y: h - r})
	path.Cube(f32.Point{X: 0, Y: r * c}, f32.Point{X: +r - r*c, Y: r}, f32.Point{X: +r, Y: r})      // SW
	path.Cube(f32.Point{X: +r * c, Y: 0}, f32.Point{X: +r, Y: -r + r*c}, f32.Point{X: +r, Y: -r})   // SE
	path.Cube(f32.Point{X: 0, Y: -r * c}, f32.Point{X: -(r - r*c), Y: -r}, f32.Point{X: -r, Y: -r}) // NE
	path.Cube(f32.Point{X: -r * c, Y: 0}, f32.Point{X: -r, Y: r - r*c}, f32.Point{X: -r, Y: r})     // NW
	path.Close()
	//clip.Rect{Max: image.Point{X: int(p.X + w), Y: int(p.Y + h)}}.Op()
	paint.ColorOp{Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
