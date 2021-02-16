package colorpicker

// https://bgrins.github.io/spectrum/#why

import (
	"fmt"
	"gioui.org/f32"
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

func NewColorPicker(th *material.Theme) *ColorPicker {
	cp := &ColorPicker{
		tone:      &Position{},
		hue:       &widget.Float{Axis: layout.Horizontal},
		alfa:      &widget.Float{Axis: layout.Horizontal},
		input:     &widget.Editor{Alignment: text.Middle, SingleLine: true},
		toggle:    &widget.Clickable{},
		hexEditor: &HexEditor{th, &widget.Editor{SingleLine: true}},
		rgbEditor: &RgbEditor{th, &widget.Editor{SingleLine: true}, &widget.Editor{SingleLine: true}, &widget.Editor{SingleLine: true}},
		hsvEditor: &HsvEditor{th, &widget.Editor{SingleLine: true}, &widget.Editor{SingleLine: true}, &widget.Editor{SingleLine: true}},
		theme:     th}
	cp.SetColor(color.RGBA{R: 255, A: 255})
	return cp
}

type ColorPicker struct {
	// Encode color saturation on X-axis and color value on y-axis.
	tone  *Position
	hue   *widget.Float
	alfa  *widget.Float
	input *widget.Editor

	hexEditor *HexEditor
	rgbEditor *RgbEditor
	hsvEditor *HsvEditor

	inputType int
	toggle    *widget.Clickable
	color     HSVColor
	theme     *material.Theme
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Point{X: gtx.Px(unit.Dp(210)), Y: gtx.Px(unit.Dp(200))})
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(cp.layoutGradiants),
		layout.Rigid(cp.layoutRainbow),
		layout.Rigid(cp.layoutAlpha),
		layout.Rigid(cp.layoutTextInput))
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
	macro := op.Record(gtx.Ops)
	var w layout.Widget
	switch cp.inputType {
	case 0:
		w = cp.hexEditor.Layout
	case 1:
		w = cp.rgbEditor.Layout
	default:
		w = cp.hsvEditor.Layout
	}
	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, cp.toggle, func(gtx layout.Context) layout.Dimensions {
				return toggleIcon.Layout(gtx, cp.theme.TextSize)
			})
		}),
		layout.Flexed(1, w))
	call := macro.Stop()
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, clip.Rect{Max: dims.Size}.Op())
	call.Add(gtx.Ops)
	return dims
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

func (cp *ColorPicker) NRGBA() color.NRGBA {
	return f32color.RGBAToNRGBA(cp.RGBA())
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
	for range cp.toggle.Clicks() {
		cp.inputType = int(math.Mod(float64(cp.inputType+1), 3))
	}
}
