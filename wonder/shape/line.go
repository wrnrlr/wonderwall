package shape

import (
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/rs/xid"
	"image/color"
)

type Path []f32.Point

type rect [4]f32.Point

func (r rect) hit(p f32.Point) bool {
	return pointInTriangle(p, r[0], r[1], r[2]) || pointInTriangle(p, r[0], r[3], r[2])
}

type Line struct {
	ID     string
	Points Path
	Color  color.NRGBA
	Width  float32
	offset f32.Point
	rects  []rect
	boxes  []f32.Rectangle
}

func NewPolyline(points []f32.Point, col color.NRGBA, width float32) *Line {
	return &Line{
		ID:     xid.New().String(),
		Points: points,
		Color:  col,
		Width:  width,
	}
}

func (l *Line) Bounds() f32.Rectangle {
	r := l.Width
	if l.boxes == nil {
		length := len(l.Points)
		for i, p1 := range l.Points {
			b := f32.Rect(p1.X-r, p1.Y-r, p1.X+r, p1.Y+r)
			l.boxes = append(l.boxes, b)
			if i < length-1 {
				p2 := l.Points[i+1]
				tilt := angle(p1, p2) + rad90
				a := offsetPoint(p1, l.Width, tilt)
				b := offsetPoint(p2, l.Width, tilt)
				c := offsetPoint(p2, -l.Width, tilt)
				d := offsetPoint(p1, -l.Width, tilt)
				l.rects = append(l.rects, rect{a, b, c, d})
				box := boundingBox([]f32.Point{a, b, c, d})
				l.boxes = append(l.boxes, box)
			}
		}
	}
	if len(l.boxes) == 0 {
		return f32.Rectangle{}
	}
	box := l.boxes[0]
	for _, b := range l.boxes[1:] {
		box = box.Union(b)
	}
	return box
}

func (l *Line) Offset(p f32.Point) Shape {
	l.offset = p
	return l
}

func (l Line) Draw(gtx C) {
	scale := gtx.Metric.PxPerDp
	width := l.Width * scale
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	for i, p := range l.Points {
		if i == 0 {
			path.MoveTo(p.Mul(scale))
		} else {
			path.LineTo(p.Mul(scale))
		}
	}
	style := clip.StrokeStyle{Width: width, Miter: 10, Cap: clip.RoundCap, Join: clip.RoundJoin}
	paint.FillShape(gtx.Ops, l.Color, clip.Stroke{Path: path.End(), Style: style}.Op())
}

func (l *Line) Move(delta f32.Point) {
	for i, p := range l.Points {
		l.Points[i] = p.Add(delta)
	}
	l.boxes = nil
	l.rects = nil
}

func (l *Line) Hit(p f32.Point) bool {
	for _, r := range l.rects {
		if r.hit(p) {
			return true
		}
	}
	return false
}

func (l *Line) Eq(s2 Shape) bool {
	return false
}

func (l *Line) Identity() string {
	return l.ID
}
