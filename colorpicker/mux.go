package colorpicker

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"image"
	"image/color"
)

type Mux struct {
	inputs []ColorInput
	color  color.NRGBA
}

func (m *Mux) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Point{X: gtx.Px(unit.Dp(210)), Y: gtx.Px(unit.Dp(200))})
	var children []layout.FlexChild
	for _, input := range m.inputs {
		children = append(children, layout.Rigid(input.Layout))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (m *Mux) Color() color.NRGBA {
	return m.color
}

func (m *Mux) SetColor(col color.NRGBA) {
	m.color = col
	for _, input := range m.inputs {
		input.SetColor(col)
	}
}

func (m *Mux) Changed() bool {
	changed := false
	for _, input := range m.inputs {
		if input.Changed() {
			changed = true
			m.color = input.Color()
		}
	}
	return changed
}
