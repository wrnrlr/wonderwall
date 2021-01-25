package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/wrnrlr/wonderwall/wonder/ui"
)

type Topbar struct {
	Back bool

	backClick *widget.Clickable
}

func NewTopbar(back bool) *Topbar {
	return &Topbar{
		Back:      back,
		backClick: new(widget.Clickable),
	}
}

func (t *Topbar) Event(gtx C) interface{} {
	if t.backClick.Clicked() {
		return BackEvent{}
	}
	return nil
}

func (t *Topbar) Layout(gtx C, insets layout.Inset, w layout.Widget) D {
	return layout.Stack{Alignment: layout.SW}.Layout(gtx,
		layout.Expanded(fill{theme.Color.Primary}.Layout),
		layout.Stacked(func(gtx C) D {
			return insets.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if !t.Back {
							return layout.Dimensions{}
						}
						return ui.Item(theme, t.backClick, backIcon).Layout(gtx)
					}),
					layout.Flexed(1, w),
				)
			})
		}),
	)
}
