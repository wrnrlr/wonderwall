package main

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type Text struct {
	f32.Point
	ID          string
	Text        string
	StrokeColor color.RGBA
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

func (t *TextWriter) Event(gtx C) (pe *pointer.Event) {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(t) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
		case pointer.Drag:
		case pointer.Release, pointer.Cancel:
			pe = &e
		}
	}
	pointer.InputOp{Tag: t, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	return pe
}
