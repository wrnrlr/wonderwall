package shape

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

type Path []f32.Point

type Polyline struct {
	Points Path
	Color  color.RGBA
	Width  float32

	boxes []image.Rectangle
}

func (l Polyline) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test

func (l Polyline) Offset(p f32.Point) Shape {
	return nil
}

func (l Polyline) Draw(ops *op.Ops) {
	l.drawPolyline(l.Points, l.Width, l.Color, ops)
}

func (l Polyline) drawPolyline(points []f32.Point, width float32, col color.RGBA, ops *op.Ops) {
	length := len(points)
	for i, p := range points {
		l.drawCircle(p, width, col, ops)
		if i < length-1 {
			l.drawLine(p, points[i+1], width, col, ops)
		}
	}
}

func (l Polyline) drawCircle(p f32.Point, radius float32, col color.RGBA, ops *op.Ops) {
	d := radius * 2
	const k = 0.551915024494 // 4*(sqrt(2)-1)/3
	defer op.Push(ops).Pop()
	var path clip.Path
	path.Begin(ops)
	path.Move(f32.Point{X: p.X + radius, Y: p.Y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.End().Add(ops)
	box := f32.Rectangle{Min: f32.Point{X: p.X - radius, Y: p.Y - radius}, Max: f32.Point{X: p.X + d, Y: p.Y + d}}
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{Rect: box}.Add(ops)
}

func (l Polyline) drawLine(p1, p2 f32.Point, width float32, col color.RGBA, gtx layout.Context) {
	tilt := angle(p1, p2)
	a := offsetPoint(p1, width, tilt+rad90)
	b := offsetPoint(p2, width, tilt+rad90)
	c := offsetPoint(p2, -width, tilt+rad90)
	d := offsetPoint(p1, -width, tilt+rad90)
	stack := op.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(a)
	path.Line(b.Sub(a))
	path.Line(c.Sub(b))
	path.Line(d.Sub(c))
	path.Line(a.Sub(d))
	path.End().Add(gtx.Ops)
	box := boundingBox([]f32.Point{a, b, c, d})
	paint.PaintOp{Rect: box}.Add(gtx.Ops)
	stack.Pop()
}

func (l Polyline) Hit(p f32.Point) bool {
	//points := l.Points
	//width := l.Width
	//length := len(points)
	//stack := op.Push(gtx.Ops)
	//for i, p1 := range points {
	//	if i < length-1 {
	//		p2 := points[i+1]
	//		tilt := angle(p1, p2)
	//		a := offsetPoint(p1, width, tilt+rad90)
	//		b := offsetPoint(p2, width, tilt+rad90)
	//		c := offsetPoint(p2, -width, tilt+rad90)
	//		d := offsetPoint(p1, -width, tilt+rad90)
	//		box := boundingBox([]f32.Point{a, b, c, d})
	//		pointer.Rect(imageRect(box)).Add(gtx.Ops)
	//		green := &color.RGBA{0, 255, 0, 255}
	//		Rectangle{box, nil, green, 2}.Draw(gtx.Ops)
	//		//pointer.InputOp{Tag: &l, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	//	}
	//}
	//for _, e := range gtx.Events(&l) {
	//	e, ok := e.(pointer.Event)
	//	if !ok {
	//		continue
	//	}
	//	fmt.Println("Pointer event")
	//	switch e.Type {
	//	case pointer.Press:
	//		fmt.Print("Press")
	//	case pointer.Drag:
	//		fmt.Print("Drag")
	//	case pointer.Release, pointer.Cancel:
	//		fmt.Print("Release, Cancel")
	//	}
	//}
	//pointer.InputOp{Tag: &l, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	//stack.Pop()
	return false
}
