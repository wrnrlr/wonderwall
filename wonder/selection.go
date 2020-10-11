package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image"
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

func (s *Selection) Event(plane *shape.Plane, gtx C) interface{} {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(s) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
			sh := plane.Hits(pos)
			if sh != nil {
				s.ToggleSelection(sh)
			} else {
				s.Clear()
			}
			fmt.Printf("results: %v\n", sh)
		case pointer.Scroll:
			if e.Modifiers.Contain(key.ModCommand) || e.Modifiers.Contain(key.ModCtrl) {
				return ZoomEvent{Scroll: e.Scroll.Y, Pos: e.Position}
			} else {
				return PanEvent{Offset: e.Scroll, Pos: e.Position}
			}
		case pointer.Drag:
		case pointer.Release, pointer.Cancel:
		}
	}
	pointer.InputOp{Tag: s, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release | pointer.Scroll}.Add(gtx.Ops)
	return nil
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
