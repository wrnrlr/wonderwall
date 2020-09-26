package main

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/Almanax/wonderwall/wonder/ui"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

type Topbar struct {
	Back bool

	backClick gesture.Click
}

func (t *Topbar) Event(gtx layout.Context) interface{} {
	for _, e := range t.backClick.Events(gtx) {
		if e.Type == gesture.TypeClick {
			return BackEvent{}
		}
	}
	return nil
}

func (t *Topbar) Layout(gtx layout.Context, insets layout.Inset, w layout.Widget) layout.Dimensions {
	insets = layout.Inset{
		Top:    unit.Add(gtx.Metric, insets.Top, unit.Dp(16)),
		Bottom: unit.Dp(16),
		Left:   unit.Max(gtx.Metric, insets.Left, unit.Dp(16)),
		Right:  unit.Max(gtx.Metric, insets.Right, unit.Dp(16)),
	}
	return layout.Stack{Alignment: layout.SW}.Layout(gtx,
		layout.Expanded(fill{theme.Color.Primary}.Layout),
		layout.Stacked(func(gtx C) D {
			return insets.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if !t.Back {
							return layout.Dimensions{}
						}
						ico := (&ui.Icon{Src: icons.NavigationArrowBack, Size: unit.Dp(24)}).Image(gtx.Metric, theme.Color.Text)
						ico.Add(gtx.Ops)
						paint.PaintOp{Rect: f32.Rectangle{Max: toPointF(ico.Size())}}.Add(gtx.Ops)
						dims := layout.Dimensions{Size: ico.Size()}
						dims.Size.X += gtx.Px(unit.Dp(4))
						pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
						t.backClick.Add(gtx.Ops)
						return dims
					}),
					layout.Flexed(1, w),
				)
			})
		}),
	)
}
