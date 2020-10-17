package shape

import (
	"gioui.org/f32"
	"gioui.org/op/paint"
	"github.com/rs/xid"
)

type Image struct {
	ID    string
	X, Y  float32
	Image *paint.ImageOp
}

func NewImage(x, y float32, img *paint.ImageOp) *Image {
	return &Image{
		ID:    xid.New().String(),
		X:     x,
		Y:     y,
		Image: img,
	}
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

func (i *Image) Move(delta f32.Point) {
	pos := f32.Point{i.X, i.Y}.Add(delta)
	i.X, i.Y = pos.X, pos.Y
}

func (i *Image) Eq(s Shape) bool {
	i2, ok := s.(*Image)
	if !ok {
		return false
	}
	return i.ID == i2.ID
}

func (i Image) Identity() string {
	return i.ID
}
