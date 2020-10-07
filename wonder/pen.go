package main

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image"
	"image/color"
)

type Pen struct {
	events []f32.Point
	grab   bool
}

func (p *Pen) Draw(gtx C, width float32, col color.RGBA) {
	if p.events != nil {
		shape.Polyline{Points: p.events, Width: width, Color: col}.Draw(gtx)
	}
}

func (p *Pen) Event(gtx C) []f32.Point {
	var l []f32.Point
	scale := 1 / gtx.Metric.PxPerDp
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(p) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		pos := e.Position.Mul(scale)
		switch e.Type {
		case pointer.Press:
			p.events = append(p.events, pos)
		case pointer.Drag:
			p.events = append(p.events, pos)
		case pointer.Release, pointer.Cancel:
			l = append(p.events, pos)
			p.events = nil
		}
	}
	pointer.InputOp{Tag: p, Grab: p.grab, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	return l
}
