package shape

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image"
	"image/color"
	"math"
)

type (
	D = layout.Dimensions
	C = layout.Context
)

const (
	rad90 = float32(90 * math.Pi / 180)
)

type Path []f32.Point

type Line struct {
	Points Path
	Color  color.RGBA
	Width  float32

	boxes []image.Rectangle
}

func (l Line) Layout(gtx C) {
	l.drawPolyline(l.Points, l.Width, l.Color, gtx)
}

func (l Line) drawPolyline(points []f32.Point, width float32, col color.RGBA, gtx layout.Context) {
	length := len(points)
	for i, p := range points {
		l.drawCircle(p, width, col, gtx)
		if i < length-1 {
			l.drawLine(p, points[i+1], width, col, gtx)
		}
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

func (l Line) drawLine(p1, p2 f32.Point, width float32, col color.RGBA, gtx layout.Context) {
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

func (l Line) Hit(gtx layout.Context) bool {
	points := l.Points
	width := l.Width
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
			r := Rectangle(box)
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
	pointer.InputOp{Tag: &l, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	stack.Pop()
	return false
}

func boundingBox(points []f32.Point) (box f32.Rectangle) {
	box.Min, box.Max = points[0], points[0]
	for _, p := range points {
		box.Min.X = min(box.Min.X, p.X)
		box.Min.Y = min(box.Min.Y, p.Y)
		box.Max.X = max(box.Max.X, p.X)
		box.Max.Y = max(box.Max.Y, p.Y)
	}
	return box
}

func offsetPoint(point f32.Point, distance, angle float32) f32.Point {
	x := point.X + distance*cos(angle)
	y := point.Y + distance*sin(angle)
	return f32.Point{X: x, Y: y}
}

func angle(p1, p2 f32.Point) float32 {
	return atan2(p2.Y-p1.Y, p2.X-p1.X)
}

func cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func imageRect(r f32.Rectangle) image.Rectangle {
	x1, y1 := int(r.Min.X), int(r.Min.Y)
	x2, y2 := int(r.Max.X), int(r.Max.Y)
	return image.Rect(x1, y1, x2, y2)
}
