package main

import (
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Almanax/wonderwall/wonder/shape"
	"golang.org/x/image/colornames"
)

type Selection struct {
	events    []f32.Point
	selection map[shape.Shape]bool
	prev      f32.Point
}

func NewSelection() *Selection {
	return &Selection{
		selection: map[shape.Shape]bool{},
	}
}

func (s Selection) Draw(plane *shape.Plane, gtx layout.Context) {
	tr := plane.GetTransform()
	defer op.Push(gtx.Ops).Pop()
	op.Affine(tr).Add(gtx.Ops)
	for sh, _ := range s.selection {
		b := sh.Bounds()
		shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &colornames.Lightblue, StrokeWidth: unit.Dp(3).V}.Draw(gtx)
	}
	shape.Circle{Center: s.prev, Radius: unit.Dp(5).V, FillColor: &red, StrokeColor: nil}.Fill(red, gtx)
}

func (s *Selection) Event(e pointer.Event, plane *shape.Plane, gtx C) interface{} {
	var result interface{}
	pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
	pos = plane.GetTransform().Invert().Transform(pos)
	switch e.Type {
	case pointer.Press:
		s.prev = pos
	case pointer.Drag:
		delta := pos.Sub(s.prev)
		s.prev = pos
		result = DragSelectionEvent{delta}
	case pointer.Release:
		if s.prev.X == pos.X && s.prev.Y == pos.Y {
			sh := plane.Hits(pos)
			if sh == nil {
				s.Clear()
			} else if e.Modifiers.Contain(key.ModShift) {
				s.ToggleSelection(sh)
			} else {
				s.SetSelection(sh)
			}
		} else {
			delta := pos.Sub(s.prev)
			result = MoveSelectionEvent{delta}
		}
	}
	return result
}

func (s *Selection) SetSelection(sh shape.Shape) {
	s.Clear()
	s.selection[sh] = true
}

func (s *Selection) ToggleSelection(sh shape.Shape) {
	if _, ok := s.selection[sh]; ok {
		delete(s.selection, sh)
	} else {
		s.selection[sh] = true
	}
}

func (s *Selection) Elements() (shapes []shape.Shape) {
	for e := range s.selection {
		shapes = append(shapes, e)
	}
	return shapes
}

func (s *Selection) Clear() {
	s.selection = map[shape.Shape]bool{}
}
