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
	"math"
)

func NewColorPicker(th *material.Theme) *ColorPicker {
	col := color.NRGBA{R: 255, A: 255}
	cp := &ColorPicker{
		tone:      &Position{},
		hue:       &widget.Float{Axis: layout.Horizontal},
		alfa:      &AlphaSlider{slider: widget.Float{Axis: layout.Horizontal}, color: col},
		toggle:    &widget.Clickable{},
		hexEditor: &HexEditor{theme: th, hex: newHexField(widget.Editor{SingleLine: true}, hex.EncodeToString([]byte{col.R, col.G, col.B}))},
		rgbEditor: &RgbEditor{theme: th, r: &byteField{Editor: widget.Editor{SingleLine: true}}, g: &byteField{Editor: widget.Editor{SingleLine: true}}, b: &byteField{Editor: widget.Editor{SingleLine: true}}},
		hsvEditor: &HsvEditor{theme: th, h: &degreeField{Editor: widget.Editor{SingleLine: true}}, s: &percentageField{Editor: widget.Editor{SingleLine: true}}, v: &percentageField{Editor: widget.Editor{SingleLine: true}}},
		theme:     th}
	cp.SetColor(col)
	return cp
}

type ColorPicker struct {
	// Encode hsv saturation on X-axis and hsv value on y-axis.
	tone *Position
	hue  *widget.Float
	alfa *AlphaSlider

	hexEditor *HexEditor
	rgbEditor *RgbEditor
	hsvEditor *HsvEditor

	inputType int
	toggle    *widget.Clickable
	hsv       HSVColor
	color     color.NRGBA
	theme     *material.Theme
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
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, cp.toggle, func(gtx layout.Context) layout.Dimensions {
				return toggleIcon.Layout(gtx, unit.Dp(30))
			})
		}),
		layout.Flexed(1, w))
	call := macro.Stop()
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, clip.Rect{Max: dims.Size}.Op())
	call.Add(gtx.Ops)
	return dims
}

func (cp *ColorPicker) SetColor(col color.NRGBA) {
	cp.setColor(col)
	cp.hexEditor.SetColor(col)
	cp.rgbEditor.SetColor(col)
	cp.hsvEditor.SetColor(col)
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
		cp.hexEditor.SetColor(col)
		//cp.rgbEditor.SetColor(col)
		//cp.hsvEditor.SetColor(col)
	}
	if cp.hue.Changed() {
		changed = true
		cp.hsv.H = cp.hue.Value * 360
		col := cp.Color()
		cp.hexEditor.SetColor(col)
		//cp.rgbEditor.SetColor(col)
		//cp.hsvEditor.SetColor(col)
	}
	if cp.alfa.Changed() {
		col := cp.Color()
		cp.hexEditor.SetColor(col)
		//cp.rgbEditor.SetColor(col)
		//cp.hsvEditor.SetColor(col)
	}
	if cp.inputType == 0 && cp.hexEditor.Changed() {
		changed = true
		col := cp.hexEditor.Color()
		cp.setColor(col)
		//cp.rgbEditor.SetColor(col)
		//cp.hsvEditor.SetColor(col)
	}
	//if cp.inputType == 1 && cp.rgbEditor.Changed() {
	//	changed = true
	//	col := cp.hexEditor.Color()
	//	cp.setColor(col)
	//	cp.hexEditor.SetColor(col)
	//	cp.hsvEditor.SetColor(col)
	//}
	//if cp.inputType == 2 && cp.hsvEditor.Changed() {
	//	changed = true
	//	col := cp.hexEditor.Color()
	//	cp.setColor(col)
	//	cp.hexEditor.SetColor(col)
	//	cp.rgbEditor.SetColor(col)
	//}
	for range cp.toggle.Clicks() {
		cp.inputType = int(math.Mod(float64(cp.inputType+1), 3))
	}
	return changed
}
