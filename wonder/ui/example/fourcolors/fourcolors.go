package main

// https://bgrins.github.io/spectrum/#why

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/wrnrlr/wonderwall/wonder/f32color"
	"image"
	"image/color"
	"math"
)

var th *material.Theme

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		th = material.NewTheme(gofont.Collection())
		colorPicker := NewColorPicker()
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

func NewColorPicker() *ColorPicker {
	cp := &ColorPicker{
		tone:  &Position{},
		hue:   &widget.Float{Axis: layout.Horizontal},
		alfa:  &widget.Float{Axis: layout.Horizontal},
		input: &widget.Editor{Alignment: text.Middle, SingleLine: true}}
	cp.SetColor(color.RGBA{R: 255, A: 255})
	return cp
}

type ColorPicker struct {
	// Encode color saturation on X-axis and color value on y-axis.
	tone  *Position
	hue   *widget.Float
	alfa  *widget.Float
	input *widget.Editor

	color HSVColor
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(image.Point{X: gtx.Px(unit.Dp(210)), Y: gtx.Px(unit.Dp(200))})
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(cp.layoutGradiants),
			layout.Rigid(cp.layoutRainbow),
			layout.Rigid(cp.layoutAlpha),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return cp.layoutTextInput(gtx)
			}))
	})
}

func (cp *ColorPicker) layoutGradiants(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Px(unit.Dp(120))
	dr := image.Rectangle{Max: image.Point{X: w, Y: h}}
	primary := f32color.RGBAToNRGBA(HsvToRgb(HSVColor{H: cp.hue.Value * 360, S: 1, V: 1}))
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
	cp.tone.Layout(gtx, 1, f32.Point{}, f32.Point{X: 1, Y: 1})
	p := cp.tone.Pos()
	drawControl(p, 10, 1, gtx)

	return layout.Dimensions{Size: dr.Max}
}

var (
	red     = color.RGBA{R: 255, A: 255}
	yellow  = color.RGBA{R: 255, G: 255, A: 255}
	green   = color.RGBA{G: 255, A: 255}
	cyan    = color.RGBA{G: 255, B: 255, A: 255}
	blue    = color.RGBA{B: 255, A: 255}
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

	gtx.Constraints = layout.Exact(image.Point{X: w, Y: h})
	cp.hue.Layout(gtx, 1, 0, 1)
	x := cp.hue.Pos()
	drawControl(f32.Point{x, float32(h / 2)}, 10, 1, gtx)

	return layout.Dimensions{Size: image.Point{X: w, Y: h}}
}

func (cp *ColorPicker) layoutAlpha(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Px(unit.Dp(20))

	gtx.Constraints = layout.Exact(image.Point{w, h})
	drawCheckerboard(gtx)

	col1 := cp.RGBA()
	col2 := col1
	col1.A = 0x00
	col2.A = 0xff
	defer op.Save(gtx.Ops).Load()
	paint.LinearGradientOp{
		Stop1:  f32.Point{float32(0), 0},
		Stop2:  f32.Point{float32(w), 0},
		Color1: f32color.RGBAToNRGBA(col1),
		Color2: f32color.RGBAToNRGBA(col2),
	}.Add(gtx.Ops)
	dr := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}}
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	//gtx.Constraints = layout.Exact(image.Point{X: w, Y: h})
	cp.alfa.Layout(gtx, 1, 0, 1)
	x := cp.alfa.Pos()
	drawControl(f32.Point{x, float32(h / 2)}, 10, 1, gtx)

	return layout.Dimensions{Size: image.Point{X: w, Y: h}}
}

func (cp *ColorPicker) layoutTextInput(gtx layout.Context) layout.Dimensions {
	es := material.Editor(th, cp.input, "hex")
	es.Font = text.Font{Variant: "Mono"}
	return es.Layout(gtx)
}

func (cp *ColorPicker) SetColor(rgb color.RGBA) {
	hsv := RgbToHsv(rgb)
	cp.tone.X = hsv.S
	cp.tone.Y = 1 - hsv.V
	cp.hue.Value = hsv.H
	cp.alfa.Value = float32(rgb.A / 255)
	cp.setText()
}

func (cp *ColorPicker) RGBA() color.RGBA {
	fmt.Printf("%v, %v, %v\n", cp.hue.Value, cp.tone.Y, cp.tone.X)
	return HsvToRgb(HSVColor{H: cp.hue.Value * 360, S: cp.tone.X, V: 1 - cp.tone.Y})
}

func (cp *ColorPicker) setText() {
	rgba := cp.RGBA()
	r, g, b, a := rgba.RGBA()
	cp.input.SetText(fmt.Sprintf("%x%x%x%x", int(r), int(g), int(b), int(a)))
}

func (cp *ColorPicker) Event() {
	if cp.tone.Changed() {
		cp.setText()
	}
	if cp.hue.Changed() {
		cp.setText()
	}
	if cp.alfa.Changed() {
		cp.setText()
	}
	cp.input.Events()
}

const c = 0.55228475 // 4*(sqrt(2)-1)/3

func drawControl(p f32.Point, radius, width float32, gtx layout.Context) {
	width = float32(gtx.Px(unit.Dp(width))) / 2
	radius = float32(gtx.Px(unit.Dp(radius))) - width
	p.X -= radius - width*4
	p.Y -= radius - width*9
	drawCircle(p, radius, width, color.NRGBA{A: 0xff}, gtx)
	p.X += width * 2
	p.Y += width * 2
	drawCircle(p, radius, width, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, gtx)
}

func drawCircle(p f32.Point, radius, width float32, col color.NRGBA, gtx layout.Context) {
	w, h := radius*2, radius*2
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Point{X: p.X, Y: p.Y})
	path.Move(f32.Point{X: w / 4 * 3, Y: radius / 2})
	path.Cube(f32.Point{X: 0, Y: radius * c}, f32.Point{X: -radius + radius*c, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * c, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*c}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * c}, f32.Point{X: radius - radius*c, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * c, Y: 0}, f32.Point{X: radius, Y: radius - radius*c}, f32.Point{X: radius, Y: radius})      // NE
	path.Move(f32.Point{X: -w, Y: -radius})                                                                                     // Return to origin
	scale := (radius - width*2) / radius
	path.Move(f32.Point{X: w * (1 - scale) * .5, Y: h * (1 - scale) * .5})
	w *= scale
	h *= scale
	radius *= scale
	path.Move(f32.Point{X: 0, Y: h - radius})
	path.Cube(f32.Point{X: 0, Y: radius * c}, f32.Point{X: +radius - radius*c, Y: radius}, f32.Point{X: +radius, Y: radius})      // SW
	path.Cube(f32.Point{X: +radius * c, Y: 0}, f32.Point{X: +radius, Y: -radius + radius*c}, f32.Point{X: +radius, Y: -radius})   // SE
	path.Cube(f32.Point{X: 0, Y: -radius * c}, f32.Point{X: -(radius - radius*c), Y: -radius}, f32.Point{X: -radius, Y: -radius}) // NE
	path.Cube(f32.Point{X: -radius * c, Y: 0}, f32.Point{X: -radius, Y: radius - radius*c}, f32.Point{X: -radius, Y: radius})     // NW
	clip.Outline{Path: path.End()}.Op().Add(gtx.Ops)
	cons := gtx.Constraints
	dr := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: cons.Max.X, Y: cons.Max.Y}}
	clip.Rect(dr).Add(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawCheckerboard(gtx layout.Context) {
	w := gtx.Constraints.Max.X
	h := gtx.Constraints.Max.Y
	size := h / 2
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	count := int(math.Ceil(float64(w / size)))
	for i := 0; i < count; i++ {
		offset := 0
		if math.Mod(float64(i), 2) == 0 {
			offset += size
		}
		path.MoveTo(f32.Point{X: float32(i * size), Y: float32(offset)})
		path.Line(f32.Point{X: float32(size)})
		path.Line(f32.Point{Y: float32(size)})
		path.Line(f32.Point{X: float32(-size)})
		path.Line(f32.Point{Y: float32(-size)})
	}
	clip.Outline{Path: path.End()}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0xef, G: 0xef, B: 0xef, A: 0xff}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
