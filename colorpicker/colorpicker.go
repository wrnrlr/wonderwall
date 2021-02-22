package colorpicker

// https://bgrins.github.io/spectrum/#why

import (
	"encoding/hex"
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/wrnrlr/wonderwall/wonder/f32color"
	"image"
	"image/color"
)

func NewColorPicker(th *material.Theme) *ColorPicker {
	col := color.NRGBA{R: 255, A: 255}
	cp := &ColorPicker{
		tone: &Position{},
		hue:  &widget.Float{Axis: layout.Horizontal},
		alfa: &AlphaSlider{slider: widget.Float{Axis: layout.Horizontal}, color: col},
		editor: NewToggle(&widget.Clickable{},
			&HexEditor{theme: th, hex: newHexField(widget.Editor{SingleLine: true}, hex.EncodeToString([]byte{col.R, col.G, col.B}))},
			&RgbEditor{theme: th, r: &byteField{Editor: widget.Editor{SingleLine: true}}, g: &byteField{Editor: widget.Editor{SingleLine: true}}, b: &byteField{Editor: widget.Editor{SingleLine: true}}},
			&HsvEditor{theme: th, h: &degreeField{Editor: widget.Editor{SingleLine: true}}, s: &percentageField{Editor: widget.Editor{SingleLine: true}}, v: &percentageField{Editor: widget.Editor{SingleLine: true}}}),
		theme: th}
	cp.SetColor(col)
	return cp
}

type ColorPicker struct {
	// Encode hsv saturation on X-axis and hsv value on y-axis.
	tone *Position
	hue  *widget.Float
	alfa *AlphaSlider

	editor *Toggle

	hsv   HSVColor
	color color.NRGBA
	theme *material.Theme
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Point{X: gtx.Px(unit.Dp(210)), Y: gtx.Px(unit.Dp(200))})
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(cp.layoutGradiants),
		layout.Rigid(cp.layoutRainbow),
		layout.Rigid(cp.alfa.Layout),
		layout.Rigid(cp.layoutTextInput))
}

func (cp *ColorPicker) layoutGradiants(gtx layout.Context) layout.Dimensions {
	w := gtx.Constraints.Max.X
	h := gtx.Px(unit.Dp(120))
	dr := image.Rectangle{Max: image.Point{X: w, Y: h}}
	primary := HsvToRgb(HSVColor{H: cp.hue.Value * 360, S: 1, V: 1})
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

func (cp *ColorPicker) layoutTextInput(gtx layout.Context) layout.Dimensions {
	return cp.editor.Layout(gtx)
}

func (cp *ColorPicker) SetColor(col color.NRGBA) {
	cp.setColor(col)
	cp.editor.SetColor(col)
}

func (cp *ColorPicker) setColor(rgb color.NRGBA) {
	hsv := RgbToHsv(rgb)
	cp.tone.X = hsv.S
	cp.tone.Y = 1 - hsv.V
	cp.hue.Value = hsv.H
	//cp.alfa.Value = float32(rgb.A / 255)
}

func (cp *ColorPicker) Color() color.NRGBA {
	fmt.Printf("%v, %v, %v\n", cp.hue.Value, cp.tone.Y, cp.tone.X)
	rgb := HsvToRgb(HSVColor{H: cp.hue.Value * 360, S: cp.tone.X, V: 1 - cp.tone.Y})
	//rgb.A = byte(cp.alfa.Value * 255)
	return rgb
}

func (cp *ColorPicker) Changed() bool {
	changed := false
	if cp.tone.Changed() {
		cp.hsv.S = cp.tone.X
		cp.hsv.V = 1 - cp.tone.Y
		col := cp.Color()
		cp.editor.SetColor(col)
	}
	if cp.hue.Changed() {
		changed = true
		cp.hsv.H = cp.hue.Value * 360
		col := cp.Color()
		cp.editor.SetColor(col)
	}
	if cp.alfa.Changed() {
		col := cp.Color()
		cp.editor.SetColor(col)
	}
	if cp.editor.Changed() {
		changed = true
		col := cp.editor.Color()
		cp.setColor(col)
	}
	return changed
}
