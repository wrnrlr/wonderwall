package colorpicker

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type EditorType int

type HexEditor struct {
	theme *material.Theme
	hex   *widget.Editor
}

func (e *HexEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "Hex").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.hex, "").Layout))
}

type RgbEditor struct {
	theme *material.Theme
	r     *widget.Editor
	g     *widget.Editor
	b     *widget.Editor
}

func (e *RgbEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "R ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.r, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "G ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.g, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "B ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.b, "").Layout))
}

type HsvEditor struct {
	theme *material.Theme
	h     *widget.Editor
	s     *widget.Editor
	v     *widget.Editor
}

func (e *HsvEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "H ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.h, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "S ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.s, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "V ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.v, "").Layout))
}
