package main

import (
	"fmt"
	"gioui.org/f32"
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

func (s *Selection) Event(plane *shape.Plane, gtx C) []f32.Point {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(s) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			sh := plane.Hits(e.Position)
			if sh != nil {
				s.ToggleSelection(sh)
			} else {
				s.Clear()
			}
			fmt.Printf("results: %v\n", sh)
		case pointer.Scroll:
			fmt.Printf("Scroll: %v, %v\n", e.Position, e.Scroll)
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

func (s *Selection) Clear() {
	s.selection = map[shape.Shape]bool{}
}
