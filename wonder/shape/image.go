package shape

import (
	"gioui.org/f32"
	"gioui.org/op/paint"
)

type Image struct {
	X, Y  float32
	Image paint.ImageOp
}

func (i Image) Bounds() f32.Rectangle {
	r := toRectF(i.Image.Rect)
	return r.Add(f32.Pt(i.X, i.Y))
}

// Hit test
func (i Image) Hit(p f32.Point) bool {
	return true
}

func (i Image) Offset(p f32.Point) Shape {
	return nil
}

func (i Image) Draw(gtx C) {
	b := i.Image.Rect
	scale := gtx.Metric.PxPerDp
	p := f32.Point{X: i.X, Y: i.Y}.Mul(scale)
	w, h := float32(b.Max.X)*scale, float32(b.Max.Y)*scale
	i.Image.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(p.X, p.Y, p.X+w, p.Y+h)}.Add(gtx.Ops)
}
