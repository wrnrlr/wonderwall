package shape

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"image/color"
)

type Rectangle f32.Rectangle

func (r Rectangle) Fill(rgba color.RGBA, gtx layout.Context) f32.Rectangle {
	p1, p2 := r.Min, r.Max
	defer op.Push(gtx.Ops).Pop()
	paint.ColorOp{rgba}.Add(gtx.Ops)
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
	return box
}

func (r Rectangle) Stroke(rgba color.RGBA, lineWidth float32, gtx layout.Context) f32.Rectangle {
	frect := f32.Rectangle(r)
	st := op.Push(gtx.Ops)
	clip.Border{
		Rect:  frect,
		Width: lineWidth,
	}.Add(gtx.Ops)
	dr := f32.Rectangle{
		Max: r.Max,
	}
	paint.ColorOp{Color: rgba}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	st.Pop()
	return dr
}
