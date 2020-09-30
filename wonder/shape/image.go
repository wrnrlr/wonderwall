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
	return toRectF(i.Image.Rect)
}

// Hit test
func (i Image) Hit(p f32.Point) bool {
	return false
}

func (i Image) Offset(p f32.Point) Shape {
	return nil
}

func (i Image) Draw(gtx C) {
	b := i.Image.Rect
	w, h := float32(b.Max.X), float32(b.Max.Y)
	ops := gtx.Ops
	i.Image.Add(ops)
	paint.PaintOp{Rect: f32.Rect(i.X, i.Y, i.X+w, i.Y+h)}.Add(ops)
}
