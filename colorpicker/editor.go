package colorpicker

import (
	"encoding/hex"
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"strconv"
	"strings"
)

type EditorType int

type HexEditor struct {
	theme *material.Theme
	color color.NRGBA
	hex   *widget.Editor
}

func (e *HexEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "Hex").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.hex, "").Layout))
}

func (e *HexEditor) Color() color.NRGBA {
	return e.color
}

func (e *HexEditor) SetColor(col color.NRGBA) {
	s := hex.EncodeToString([]byte{col.R, col.G, col.B})
	e.hex.SetText(s)
	e.Changed()
}

func (e *HexEditor) Changed() bool {
	if b, ok := hexFromEditor(e.hex); ok {
		e.color.R = b[0]
		e.color.G = b[1]
		e.color.B = b[2]
		return true
	}
	return false
}

type RgbEditor struct {
	theme *material.Theme
	rgb   color.NRGBA
	r     *widget.Editor
	g     *widget.Editor
	b     *widget.Editor
}

func (e *RgbEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "R ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.r, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "G ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.g, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "B ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.b, "").Layout))
}

func (e *RgbEditor) Changed() bool {
	changed := false
	if b, ok := byteFromEditor(e.r); ok {
		e.rgb.R = b
		changed = true
	}
	if b, ok := byteFromEditor(e.g); ok {
		e.rgb.R = b
		changed = true
	}
	if b, ok := byteFromEditor(e.b); ok {
		e.rgb.B = b
		changed = true
	}
	return changed
}

func (e *RgbEditor) SetColor(col color.NRGBA) {
	e.rgb.R = col.R
	e.rgb.G = col.G
	e.rgb.B = col.B
	e.rgb.A = col.A
	e.r.SetText(strconv.Itoa(int(col.R)))
	e.g.SetText(strconv.Itoa(int(col.G)))
	e.b.SetText(strconv.Itoa(int(col.B)))
	e.Changed()
}

type HsvEditor struct {
	theme *material.Theme
	hsv   HSVColor
	h     *widget.Editor
	s     *widget.Editor
	v     *widget.Editor
}

func (e *HsvEditor) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "H ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.h, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "S ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.s, "").Layout),
		layout.Rigid(material.Label(e.theme, unit.Sp(14), "V ").Layout),
		layout.Flexed(1, material.Editor(e.theme, e.v, "").Layout))
}

func (e *HsvEditor) Color() color.NRGBA {
	return HsvToRgb(e.hsv)
}

func (e *HsvEditor) Changed() bool {
	changed := false
	if h, ok := degreeFromEditor(e.h); ok {
		e.hsv.H = float32(h)
		changed = true
	}
	if b, ok := percentageFromEditor(e.s); ok {
		e.hsv.S = float32(b / 100)
		changed = true
	}
	if b, ok := percentageFromEditor(e.v); ok {
		e.hsv.V = float32(b / 100)
		changed = true
	}
	return changed
}

func (e *HsvEditor) SetColor(col color.NRGBA) {
	hsv := RgbToHsv(col)
	e.hsv.H = hsv.H
	e.hsv.S = hsv.S
	e.hsv.V = hsv.V
	e.h.SetText(strconv.Itoa(int(hsv.H)))
	e.s.SetText(strconv.Itoa(int(hsv.S * 100)))
	e.v.SetText(strconv.Itoa(int(hsv.V * 100)))
	e.Changed()
}

func byteFromEditor(e *widget.Editor) (byte, bool) {
	b := byte(0)
	changed := false
	for _, ev := range e.Events() {
		_, ok := ev.(widget.ChangeEvent)
		if !ok {
			continue
		}
		s := e.Text()
		i, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			continue
		}
		b, changed = byte(i), true
	}
	return b, changed
}

func percentageFromEditor(e *widget.Editor) (int, bool) {
	percentage := 0
	changed := false
	for _, ev := range e.Events() {
		_, ok := ev.(widget.ChangeEvent)
		if !ok {
			continue
		}
		s := e.Text()
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			continue
		}
		b := int(i)
		percentage, changed = b, true
	}
	return percentage, changed
}

func degreeFromEditor(e *widget.Editor) (int, bool) {
	degree := 0
	changed := false
	for _, ev := range e.Events() {
		_, ok := ev.(widget.ChangeEvent)
		if !ok {
			continue
		}
		s := e.Text()
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			continue
		}
		degree, changed = int(i), true
	}
	return degree, changed
}

func hexFromEditor(e *widget.Editor) ([]byte, bool) {
	h, changed := []byte(nil), false
	for _, ev := range e.Events() {
		ce, ok := ev.(widget.ChangeEvent)
		if !ok {
			continue
		}
		s := e.Text()
		out, err := hex.DecodeString(s)
		if err != nil {
			continue
		}
		if len(out) < 3 {
			continue
		}
		fmt.Printf("%v\n", ce)
		h, changed = out, true
	}
	return h, changed
}

func valueString(in uint8) string {
	s := strconv.Itoa(int(in))
	delta := 3 - len(s)
	if delta > 0 {
		s = strings.Repeat(" ", delta) + s
	}
	return s
}
