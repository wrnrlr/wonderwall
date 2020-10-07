package shape

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image/color"
)

type Rectangle struct {
	f32.Rectangle
	FillColor   *color.RGBA
	StrokeColor *color.RGBA
	StrokeWidth float32
}

func (r Rectangle) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

func (r Rectangle) Hit(p f32.Point) bool {
	return false
}

func (r Rectangle) Offset(p f32.Point) Shape {
	return nil
}

func (r Rectangle) Draw(gtx layout.Context) {
	scale := gtx.Metric.PxPerDp
	rec := f32.Rectangle{Min: r.Rectangle.Min.Mul(scale), Max: r.Rectangle.Max.Mul(scale)}
	if r.StrokeColor != nil {
		width := r.StrokeWidth * scale
		stack := op.Push(gtx.Ops)
		clip.Border{Rect: rec, Width: width}.Add(gtx.Ops)
		dr := f32.Rectangle{Max: rec.Max}
		paint.ColorOp{Color: *r.StrokeColor}.Add(gtx.Ops)
		paint.PaintOp{Rect: dr}.Add(gtx.Ops)
		stack.Pop()
	}
	if r.FillColor != nil {
		p1, p2 := rec.Min, rec.Max
		stack := op.Push(gtx.Ops)
		paint.ColorOp{Color: *r.FillColor}.Add(gtx.Ops)
		var path clip.Path
		path.Begin(gtx.Ops)
		path.Move(p1)
		path.Line(f32.Point{X: p2.X, Y: 0})
		path.Line(f32.Point{X: 0, Y: p2.Y})
		path.Line(f32.Point{X: -p2.X, Y: 0})
		path.Line(f32.Point{X: 0, Y: -p2.Y})
		path.End().Add(gtx.Ops)
		box := f32.Rectangle{Min: p1, Max: p2}
		paint.PaintOp{box}.Add(gtx.Ops)
		stack.Pop()
	}
}
