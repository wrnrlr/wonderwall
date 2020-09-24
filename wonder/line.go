package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/Almanax/wonderwall/wonder/daabbt"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image/color"
)

type rect struct {
	a, b, c, d f32.Point
	parent     *Line
}

func newRect(p1, p2 f32.Point, width float32, parent *Line) rect {
	tilt := angle(p1, p2) + rad90
	return rect{
		a:      offsetPoint(p1, width, tilt),
		b:      offsetPoint(p2, width, tilt),
		c:      offsetPoint(p2, -width, tilt),
		d:      offsetPoint(p1, -width, tilt),
		parent: parent}
}

func (r rect) Bounds() f32.Rectangle {
	return boundingBox([]f32.Point{r.a, r.b, r.c, r.d})
}

func (r rect) draw(col color.RGBA, gtx C) {
	stack := op.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(r.a)
	path.Line(r.b.Sub(r.a))
	path.Line(r.c.Sub(r.b))
	path.Line(r.d.Sub(r.c))
	path.Line(r.a.Sub(r.d))
	path.End().Add(gtx.Ops)
	box := boundingBox([]f32.Point{r.a, r.b, r.c, r.d})
	paint.PaintOp{Rect: box}.Add(gtx.Ops)
	stack.Pop()
}

type circle [3]f32.Point

func (r circle) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

type Line struct {
	Points []f32.Point
	Width  float32

	parts []rect
}

func (l *Line) Draw(conf *PenConfig, gtx C) {
	if l.parts == nil {
		l.makeParts()
	}
	l.drawPolyline(l.Points, l.Width, conf.StrokeColor, gtx)
}

func (l Line) drawPolyline(points []f32.Point, width float32, col color.RGBA, gtx layout.Context) {
	for _, p := range points {
		l.drawCircle(p, width, col, gtx)
	}
	for i := range l.parts {
		l.parts[i].draw(col, gtx)
	}
}

func (l Line) drawCircle(p f32.Point, radius float32, col color.RGBA, gtx layout.Context) {
	d := radius * 2
	const k = 0.551915024494 // 4*(sqrt(2)-1)/3
	defer op.Push(gtx.Ops).Pop()
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Point{X: p.X + radius, Y: p.Y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.End().Add(gtx.Ops)
	box := f32.Rectangle{Min: f32.Point{X: p.X - radius, Y: p.Y - radius}, Max: f32.Point{X: p.X + d, Y: p.Y + d}}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: box}.Add(gtx.Ops)
}

func (l *Line) makeParts() {
	length := len(l.Points)
	for i, p := range l.Points {
		if i < length-1 {
			l.parts = append(l.parts, newRect(p, l.Points[i+1], l.Width, l))
		}
	}
}

func (l *Line) Register(tree *daabbt.Node) {
	if l.parts == nil {
		l.makeParts()
	}
	for i := range l.parts {
		tree.Insert(l.parts[i])
	}
}

func (l *Line) Unregister(tree *daabbt.Node) {

}

func (l *Line) boxes(gtx C) {
	for i := range l.parts {
		box := l.parts[i].Bounds()
		stack := op.Push(gtx.Ops)
		r := shape.Rectangle(box)
		r.Stroke(color.RGBA{0, 255, 0, 125}, 2, gtx)
		stack.Pop()
	}
}

func (l *Line) Event(conf *PenConfig, gtx C) {
	points := l.Points
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
