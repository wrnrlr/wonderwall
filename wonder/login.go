package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/Almanax/wonderwall/wonder/ui"
)

type LoginPage struct {
	env     *Env
	list    *layout.List
	fields  []*Field
	submit  *widget.Clickable
	account *Account
}

func NewLoginPage(env *Env) *LoginPage {
	account := &Account{}
	p := &LoginPage{
		env:  env,
		list: &layout.List{Axis: layout.Vertical},
		fields: []*Field{
			{Header: "Email address", Hint: "you@example.org", Value: &account.User},
			{Header: "Password", Hint: "correct horse battery staple", Value: &account.Password}},
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

func (p *LoginPage) Start(stop <-chan struct{}) {}

func (p *LoginPage) Event(gtx C) interface{} {
	if p.submit.Clicked() {
		for _, f := range p.fields {
			*f.Value = f.edit.Text()
		}
		return LoginEvent{}
	}
	return nil
}

func (p *LoginPage) Layout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			var t Topbar
			return t.Layout(gtx, p.env.insets, func(gtx C) D {
				lbl := ui.H6(theme, "Sign in")
				return lbl.Layout(gtx)
			})
		}),
		layout.Flexed(1, p.layoutSigninForm),
	)
}

func (p *LoginPage) layoutSigninForm(gtx layout.Context) layout.Dimensions {
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
			return layout.E.Layout(gtx, func(gtx C) D {
				return in.Layout(gtx, ui.Button(theme, p.submit, "Sign in").Layout)
			})
		}
	})
}
