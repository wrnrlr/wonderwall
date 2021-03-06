package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/wrnrlr/wonderwall/wonder/ui"
)

type UserPage struct {
	env    *Env
	topbar *Topbar
	list   *layout.List
}

func NewUserPage(env *Env) *UserPage {
	return &UserPage{env: env, topbar: NewTopbar(true)}
}

func (p *UserPage) Start(stop <-chan struct{}) {}

func (p *UserPage) Event(gtx C) interface{} {
	if e := p.topbar.Event(gtx); e != nil {
		return e
	}
	return nil
}

func (p *UserPage) Layout(gtx C) D {
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

func (p *UserPage) LayoutMenu(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return ui.Title(theme, "Account", gtx)
		}))
}

func (p *UserPage) LayoutContent(gtx C) D {
	return D{}
}
