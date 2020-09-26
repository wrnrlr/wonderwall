package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/Almanax/wonderwall/wonder/ui"
)

type AddWallPage struct {
	env    *Env
	topbar *Topbar
	list   *layout.List
}

func NewAddWallPage(env *Env) *AddWallPage {
	return &AddWallPage{env: env, topbar: &Topbar{Back: true}}
}

func (p *AddWallPage) Start(stop <-chan struct{}) {}

func (p *AddWallPage) Event(gtx C) interface{} {
	if e := p.topbar.Event(gtx); e != nil {
		return e
	}
	return nil
}

func (p *AddWallPage) Layout(gtx C) D {
	insets := layout.Inset{Left: unit.Dp(16), Right: unit.Dp(6)}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return p.topbar.Layout(gtx, layout.Inset{}, p.LayoutMenu)
		}),
		layout.Flexed(1, func(gtx C) D {
			return insets.Layout(gtx, func(gtx C) D {
				return p.LayoutContent(gtx)
			})
		}))
}

func (p *AddWallPage) LayoutMenu(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return ui.Label(theme, unit.Dp(22), "New Wall").Layout(gtx)
		}))
}

func (p *AddWallPage) LayoutContent(gtx C) D {
	return D{}
}
