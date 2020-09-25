package main

import "gioui.org/layout"

type JoinPage struct {
	env  *Env
	list *layout.List
}

func NewJoinPage(env *Env) *JoinPage {
	return &JoinPage{
		env:  env,
		list: &layout.List{Axis: layout.Vertical}}
}

func (p *JoinPage) Start(stop <-chan struct{}) {}
func (p *JoinPage) Event(gtx C) interface{}    { return nil }
func (p *JoinPage) Layout(gtx C) D             { return D{} }
