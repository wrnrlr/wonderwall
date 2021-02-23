package colorpicker

import (
	"encoding/hex"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"strconv"
)

func NewHexEditor(th *material.Theme) *HexEditor {
	return &HexEditor{theme: th, hex: newHexField(th, widget.Editor{SingleLine: true})}
}

type HexEditor struct {
	theme *material.Theme
	color color.NRGBA
	hex   *hexField
}

func (e *HexEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "Hex ").Layout),
		layout.Flexed(1, e.hex.Layout))
}

func (e *HexEditor) Color() color.NRGBA {
	return e.color
}

func (e *HexEditor) SetColor(col color.NRGBA) {
	e.color = col
	e.hex.SetHex([]byte{col.R, col.G, col.B})
}

func (e *HexEditor) Changed() bool {
	return e.hex.Changed()
}

func NewRgbEditor(th *material.Theme) *RgbEditor {
	return &RgbEditor{theme: th, r: &byteField{Editor: widget.Editor{SingleLine: true}}, g: &byteField{Editor: widget.Editor{SingleLine: true}}, b: &byteField{Editor: widget.Editor{SingleLine: true}}}
}

type RgbEditor struct {
	theme *material.Theme
	rgb   color.NRGBA
	r     *byteField
	g     *byteField
	b     *byteField
}

func (e *RgbEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "R ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.r.Editor, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "G ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.g.Editor, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "B ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.b.Editor, "").Layout))
}

func (e *RgbEditor) Changed() bool {
	return e.r.Changed() || e.g.Changed() || e.b.Changed()
}

func (e *RgbEditor) SetColor(col color.NRGBA) {
	e.rgb = col
	e.rgb.R = col.R
	e.rgb.G = col.G
	e.rgb.B = col.B
	e.rgb.A = col.A
	e.r.SetText(strconv.Itoa(int(col.R)))
	e.g.SetText(strconv.Itoa(int(col.G)))
	e.b.SetText(strconv.Itoa(int(col.B)))
}

func (e *RgbEditor) Color() color.NRGBA {
	return e.rgb
}

func NewHsvEditor(th *material.Theme) *HsvEditor {
	return &HsvEditor{theme: th, h: &degreeField{Editor: widget.Editor{SingleLine: true}}, s: &percentageField{Editor: widget.Editor{SingleLine: true}}, v: &percentageField{Editor: widget.Editor{SingleLine: true}}}
}

type HsvEditor struct {
	theme *material.Theme
	hsv   HSVColor
	h     *degreeField
	s     *percentageField
	v     *percentageField
}

func (e *HsvEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "H ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.h.Editor, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "S ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.s.Editor, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "V ").Layout),
		layout.Flexed(1, material.Editor(e.theme, &e.v.Editor, "").Layout))
}

func (e *HsvEditor) Color() color.NRGBA {
	return HsvToRgb(e.hsv)
}

func (e *HsvEditor) Changed() bool {
	return e.h.Changed() || e.s.Changed() || e.v.Changed()
}

func (e *HsvEditor) SetColor(col color.NRGBA) {
	hsv := RgbToHsv(col)
	e.hsv.H = hsv.H
	e.hsv.S = hsv.S
	e.hsv.V = hsv.V
	e.h.SetText(strconv.Itoa(int(hsv.H)))
	e.s.SetText(strconv.Itoa(int(hsv.S * 100)))
	e.v.SetText(strconv.Itoa(int(hsv.V * 100)))
}

func parseHex(s string) ([]byte, bool) {
	out, err := hex.DecodeString(s)
	if err != nil {
		return nil, false
	}
	if len(out) < 3 {
		return nil, false
	}
	return out, true
}

func parseByte(s string) (byte, bool) {
	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, false
	}
	return byte(i), true
}

func parseDegree(s string) (int, bool) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return int(i), false
	}
	return 0, true
}

func parsePercentage(s string) (int, bool) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return int(i), false
	}
	return 0, true
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
