package colorpicker

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
)

func NewColorField(th *material.Theme) *ColorField {
	return &ColorField{
		CornerRadius: unit.Dp(4),
		Inset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		clicker: &widget.Clickable{},
		picker:  NewColorPicker(th),
		theme:   th,
	}
}

type ColorField struct {
	CornerRadius unit.Value
	Inset        layout.Inset
	clicker      *widget.Clickable
	picker       *ColorPicker
	active       bool
	theme        *material.Theme
}

const goldenRatio = 1.618

func (cf *ColorField) Layout(gtx layout.Context) layout.Dimensions {
	h := int(2.3 * float32(gtx.Metric.Px(cf.theme.TextSize)))
	w := int(float32(h) * goldenRatio)
	size := image.Point{X: w, Y: int(h)}
	gtx.Constraints = layout.Exact(size)
	dims := material.ButtonLayoutStyle{
		Background:   cf.picker.NRGBA(),
		CornerRadius: cf.CornerRadius,
		Button:       cf.clicker,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return cf.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})
	})
	if cf.active {
		macro := op.Record(gtx.Ops)
		op.Offset(f32.Point{Y: float32(dims.Size.Y)}).Add(gtx.Ops)
		cf.picker.Layout(gtx)
		op.Defer(gtx.Ops, macro.Stop())
	}
	return dims
}

func (cf *ColorField) Event() {
	for range cf.clicker.Clicks() {
		cf.active = !cf.active
	}
	cf.picker.Event()
}

func (cf *ColorField) Click() {
	cf.clicker.Click()
}
