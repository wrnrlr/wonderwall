package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Almanax/wonderwall/wonder/shape"
)

type Selection struct {
	events    []f32.Point
	selection map[shape.Shape]bool
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
		shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &green, StrokeWidth: unit.Dp(1).V}.Draw(gtx)
	}
}

func (s *Selection) Event(e pointer.Event, plane *shape.Plane, gtx C) interface{} {
	switch e.Type {
	case pointer.Press:
		pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
		sh := plane.Hits(pos)
		if sh == nil {
			s.Clear()
		} else if e.Modifiers.Contain(key.ModShift) {
			s.ToggleSelection(sh)
		} else {
			s.SetSelection(sh)
		}
		fmt.Printf("Selected: %v\n", sh)
	case pointer.Drag:
	case pointer.Release, pointer.Cancel:
	}
	return nil
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
