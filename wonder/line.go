package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image/color"
)

type Line struct {
	points []f32.Point
}

func (l *Line) Draw(conf *PenConfig, gtx C) {
	shape.Line{Points: l.points, Color: conf.StrokeColor, Width: float32(conf.StrokeSize)}.Layout(gtx)
}

func (l *Line) Event(conf *PenConfig, gtx C) {
	points := l.points
	width := float32(conf.StrokeSize)
	length := len(points)
	stack := op.Push(gtx.Ops)
	for i, p1 := range points {
		if i < length-1 {
			p2 := points[i+1]
			tilt := angle(p1, p2)
			a := offsetPoint(p1, width, tilt+rad90)
			b := offsetPoint(p2, width, tilt+rad90)
			c := offsetPoint(p2, -width, tilt+rad90)
			d := offsetPoint(p1, -width, tilt+rad90)
			box := boundingBox([]f32.Point{a, b, c, d})
			pointer.Rect(imageRect(box)).Add(gtx.Ops)
			r := shape.Rectangle(box)
			r.Stroke(color.RGBA{0, 255, 0, 125}, 2, gtx)
			//pointer.InputOp{Tag: &l, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
		}
	}
	for _, e := range gtx.Events(&l) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		fmt.Println("Pointer event")
		switch e.Type {
		case pointer.Press:
			fmt.Print("Press")
		case pointer.Drag:
			fmt.Print("Drag")
		case pointer.Release, pointer.Cancel:
			fmt.Print("Release, Cancel")
		}
	}
	pointer.InputOp{Tag: l, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	stack.Pop()
}
