package main

import (
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"github.com/Almanax/wonderwall/wonder/shape"
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

func (p *Pen) Event(e pointer.Event, gtx C) interface{} {
	var result interface{}
	pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
	switch e.Type {
	case pointer.Press:
		p.events = []f32.Point{pos}
	case pointer.Drag:
		p.events = append(p.events, pos)
	case pointer.Release, pointer.Cancel:
		p.events = append(p.events, pos)
		result = AddLineEvent{Points: p.events}
		p.events = nil
	case pointer.Scroll:
		if e.Modifiers.Contain(key.ModCommand) || e.Modifiers.Contain(key.ModCtrl) {
			result = ZoomEvent{Scroll: e.Scroll.Y, Pos: pos}
		} else {
			result = PanEvent{Offset: e.Scroll, Pos: pos}
		}
	}
	return result
}
