package main

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/Almanax/wonderwall/wonder/ui"
)

type Field struct {
	env    *Env
	Header string
	Hint   string
	Value  *string
	edit   *widget.Editor
}

func (f *Field) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx ui.C) ui.D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			header := ui.Caption(theme, f.Header)
			header.Font.Weight = text.Bold
			dims := header.Layout(gtx)
			dims.Size.Y += gtx.Px(unit.Dp(4))
			return dims
		}),
		layout.Rigid(ui.Editor(theme, f.edit, f.Hint).Layout),
	)
}
