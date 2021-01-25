package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/wrnrlr/wonderwall/wonder/ui"
)

type AddWallPage struct {
	env    *Env
	list   *layout.List
	topbar *Topbar
	fields []*Field
	submit *widget.Clickable
}

type wallForm struct {
	Workspace string
	Title     string
}

func NewAddWallPage(env *Env) *AddWallPage {
	form := &wallForm{}
	p := &AddWallPage{
		env:    env,
		list:   &layout.List{Axis: layout.Vertical},
		topbar: NewTopbar(true),
		fields: []*Field{
			{Header: "Workspace", Hint: "", Value: &form.Workspace},
			{Header: "Title", Hint: "", Value: &form.Title}},
		submit: &widget.Clickable{}}
	for _, f := range p.fields {
		f.env = p.env
		f.edit = &widget.Editor{
			SingleLine: true,
		}
		f.edit.SetText(*f.Value)
	}
	return p
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
			return ui.Title(theme, "New Wall", gtx)
		}))
}

func (p *AddWallPage) LayoutContent(gtx C) D {
	l := p.list
	inset := layout.Inset{
		Left:  unit.Max(gtx.Metric, unit.Dp(32), p.env.insets.Left),
		Right: unit.Max(gtx.Metric, unit.Dp(32), p.env.insets.Right),
	}
	return l.Layout(gtx, len(p.fields)+1, func(gtx C, i int) D {
		in := inset
		switch {
		case i < len(p.fields):
			in.Bottom = unit.Dp(12)
			if i == 0 {
				in.Top = unit.Dp(32)
			}
			return in.Layout(gtx, p.fields[i].Layout)
		default:
			in.Bottom = unit.Max(gtx.Metric, unit.Dp(32), p.env.insets.Bottom)
			return layout.W.Layout(gtx, func(gtx C) D {
				return in.Layout(gtx, ui.Button(theme, p.submit, "Create").Layout)
			})
		}
	})
}
