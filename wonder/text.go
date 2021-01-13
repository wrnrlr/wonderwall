package main

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
)

type Text struct {
	ID          string
	X, Y        float32
	Text        string
	StrokeColor color.NRGBA
	FontWidth   float32
}

func (t Text) Draw(gtx C) {
	defer op.Push(gtx.Ops).Pop()
	op.Offset(f32.Point{X: t.X, Y: t.Y - t.FontWidth}).Add(gtx.Ops)
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(t.FontWidth), t.Text)
	l.Color = t.StrokeColor
	l.Layout(gtx)
}

type TextWriter struct {
	pointer gesture.Click
}

func (t *TextWriter) Event(e pointer.Event, gtx C) interface{} {
	var result interface{}
	pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
	switch e.Type {
	case pointer.Release, pointer.Cancel:
		result = AddTextEvent{Position: pos}
	case pointer.Scroll:
		if e.Modifiers.Contain(key.ModCommand) || e.Modifiers.Contain(key.ModCtrl) {
			result = ZoomEvent{Scroll: e.Scroll.Y, Pos: pos}
		} else {
			result = PanEvent{Offset: e.Scroll, Pos: pos}
		}
	}
	return result
}
