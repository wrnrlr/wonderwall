package shape

import (
	"gioui.org/f32"
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

func (r Rectangle) Draw(ops *op.Ops) {
	if r.StrokeColor != nil {
		stack := op.Push(ops)
		clip.Border{Rect: r.Rectangle, Width: r.StrokeWidth}.Add(ops)
		dr := f32.Rectangle{Max: r.Max}
		paint.ColorOp{Color: *r.StrokeColor}.Add(ops)
		paint.PaintOp{Rect: dr}.Add(ops)
		stack.Pop()
	}
	if r.FillColor != nil {
		p1, p2 := r.Min, r.Max
		stack := op.Push(ops)
		paint.ColorOp{Color: *r.FillColor}.Add(ops)
		var path clip.Path
		path.Begin(ops)
		path.Move(p1)
		path.Line(f32.Point{X: p2.X, Y: 0})
		path.Line(f32.Point{X: 0, Y: p2.Y})
		path.Line(f32.Point{X: -p2.X, Y: 0})
		path.Line(f32.Point{X: 0, Y: -p2.Y})
		path.End().Add(ops)
		box := f32.Rectangle{Min: p1, Max: p2}
		paint.PaintOp{box}.Add(ops)
		stack.Pop()
	}
}
