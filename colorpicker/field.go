package colorpicker

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func NewColorField(th *material.Theme) *ColorField {
	return &ColorField{
		clicker: &widget.Clickable{},
		picker:  NewColorPicker(th),
		theme:   th,
	}
}

type ColorField struct {
	clicker *widget.Clickable
	picker  *ColorPicker
	active  bool
	theme   *material.Theme
}

func (cf *ColorField) Layout(gtx layout.Context) layout.Dimensions {
	if cf.active {
		call := op.Record(gtx.Ops)
		cf.picker.Layout(gtx)
		op.Defer(gtx.Ops, call.Stop())
	}
	return material.Clickable(gtx, cf.clicker, func(gtx layout.Context) layout.Dimensions {
		cons := gtx.Constraints
		paint.FillShape(gtx.Ops, cf.picker.NRGBA(), clip.Rect{Max: cons.Max}.Op())
		return layout.Dimensions{Size: cons.Max}
	})
}

func (cf *ColorField) Event() {
	for range cf.clicker.Clicks() {
		cf.active = !cf.active
	}
	cf.picker.Event()
}
