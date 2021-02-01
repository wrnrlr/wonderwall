package main

// https://bgrins.github.io/spectrum/#why

import (
	"fmt"
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
		colorPicker := &ColorPicker{square: &Position{}}
		//white := f32color.RGBAToNRGBA(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: quarter})
		//black := f32color.RGBAToNRGBA(color.RGBA{A: quarter})
		//transparent := f32color.RGBAToNRGBA(color.RGBA{A: 0xff})
		var ops op.Ops
		for {
			e := <-w.Events()
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				colorPicker.Event()
				colorPicker.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}

type ColorPicker struct {
	square  *Position
	rainbow widget.Float
	alfa    widget.Float
	input   widget.Editor
}

var primary = f32color.RGBAToNRGBA(color.RGBA{G: 0xff, A: 0xff})

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(image.Point{X: gtx.Px(unit.Dp(210)), Y: gtx.Px(unit.Dp(200))})
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(cp.layoutGradiants),
			layout.Rigid(cp.layoutRainbow),
			layout.Rigid(cp.layoutAlpha))
	})
}

func (cp *ColorPicker) layoutGradiants(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
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
	stack.Load()

	gtx.Constraints = layout.Exact(image.Point{X: w, Y: h})
	cp.square.Layout(gtx, 1, f32.Point{}, f32.Point{X: 1, Y: 1})
	p := cp.square.Pos()
	fmt.Printf("Pos: %v\n", p)

	drawCircle(p, gtx)

	return layout.Dimensions{Size: dr.Max}
}

var (
	red     = color.RGBA{R: 255, A: 255}
	yellow  = color.RGBA{R: 255, G: 255, A: 255}
	green   = color.RGBA{G: 255, A: 255}
	cyan    = color.RGBA{G: 255, B: 255, A: 255}
	skyblue = color.RGBA{G: 127, B: 255, A: 255}
	blue    = color.RGBA{B: 255, A: 255}
	purple  = color.RGBA{R: 127, G: 0, B: 255, A: 255}
	magenta = color.RGBA{R: 255, B: 255, A: 255}
)

var colors = []color.RGBA{red, yellow, green, cyan, blue, magenta, red}

func (cp *ColorPicker) layoutRainbow(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Px(unit.Dp(20))
	stepCount := len(colors)
	stepWidth := float32(w / (stepCount - 1))
	offsetX := float32(0)
	color1 := f32color.RGBAToNRGBA(colors[0])
	dr := image.Rectangle{Max: image.Point{X: int(stepWidth), Y: h}}
	for _, col := range colors[1:] {
		color2 := f32color.RGBAToNRGBA(col)
		stack := op.Save(gtx.Ops)
		paint.LinearGradientOp{
			Stop1:  f32.Point{offsetX, 0},
			Stop2:  f32.Point{offsetX + stepWidth, 0},
			Color1: color1,
			Color2: color2,
		}.Add(gtx.Ops)
		clip.Rect(dr).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Load()
		color1 = color2
		offsetX += stepWidth
		dr = image.Rectangle{Min: image.Point{X: int(offsetX), Y: 0}, Max: image.Point{X: int(offsetX + stepWidth), Y: h}}
	}
	return layout.Dimensions{Size: image.Point{X: w, Y: h}}
}

func (cp *ColorPicker) layoutAlpha(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Px(unit.Dp(20))
	defer op.Save(gtx.Ops).Load()
	paint.LinearGradientOp{
		Stop1:  f32.Point{float32(0), 0},
		Stop2:  f32.Point{float32(w), 0},
		Color1: color.NRGBA{G: 255, A: 0},
		Color2: color.NRGBA{G: 255, A: 255},
	}.Add(gtx.Ops)
	dr := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{}
}

func (cp *ColorPicker) Event() {
	if cp.rainbow.Changed() {
	}
	if cp.alfa.Changed() {
	}
}

const c = 0.55228475 // 4*(sqrt(2)-1)/3

func drawCircle(p f32.Point, gtx layout.Context) {
	width := float32(1.5)
	r := float32(15)
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
	clip.Outline{Path: path.End()}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
