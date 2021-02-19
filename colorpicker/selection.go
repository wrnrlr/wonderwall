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

func NewColorSelection(th *material.Theme) *ColorSelection {
	return &ColorSelection{
		Dropdown:     layout.SW,
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

type ColorSelection struct {
	Dropdown     layout.Direction
	CornerRadius unit.Value
	Inset        layout.Inset
	clicker      *widget.Clickable
	picker       *ColorPicker
	active       bool
	theme        *material.Theme
}

const goldenRatio = 1.618

func (cf *ColorSelection) Layout(gtx layout.Context) layout.Dimensions {
	h := int(2.3 * float32(gtx.Metric.Px(cf.theme.TextSize)))
	w := int(float32(h) * goldenRatio)
	size := image.Point{X: w, Y: int(h)}
	gtx.Constraints = layout.Exact(size)
	dims1 := material.ButtonLayoutStyle{
		Background:   cf.picker.RGBA(),
		CornerRadius: cf.CornerRadius,
		Button:       cf.clicker,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return cf.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})
	})
	if cf.active {
		macro := op.Record(gtx.Ops)
		dims2 := cf.picker.Layout(gtx)
		var offset f32.Point
		switch cf.Dropdown {
		case layout.NE:
			offset = f32.Point{X: float32(dims1.Size.X - dims2.Size.X), Y: float32(-dims2.Size.Y)}
		case layout.SE:
			offset = f32.Point{X: float32(dims1.Size.X - dims2.Size.X), Y: float32(dims1.Size.Y)}
		case layout.SW:
			offset = f32.Point{Y: float32(dims1.Size.Y)}
		case layout.NW:
			offset = f32.Point{Y: float32(-dims2.Size.Y)}
		}
		call := macro.Stop()
		op.Offset(offset).Add(gtx.Ops)
		op.Defer(gtx.Ops, call)
	}
	return dims1
}

func (cf *ColorSelection) Event() {
	for range cf.clicker.Clicks() {
		cf.active = !cf.active
	}
	cf.picker.Changed()
}

func (cf *ColorSelection) Click() {
	cf.clicker.Click()
}
