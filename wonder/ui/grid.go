package ui

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"image"
	"math"
)

type Grid struct {
	Columns, Rows int
	Width, Height int
}

type Cell func(i, j int, gtx C) D

func (g *Grid) Layout(gtx C, cell Cell) D {
	w := gtx.Metric.Px(unit.Dp(float32(g.Width)))
	h := gtx.Metric.Px(unit.Dp(float32(g.Height)))
	columnWidth := w / g.Columns
	rowHeight := h / g.Rows
	size := image.Pt(columnWidth, rowHeight)
	gtx.Constraints = layout.Constraints{Min: size, Max: size}
	stack1 := op.Save(gtx.Ops)
	for i := 0; i < g.Rows; i++ {
		stack2 := op.Save(gtx.Ops)
		for j := 0; j < g.Columns; j++ {
			cell(i, j, gtx)
			op.Offset(f32.Pt(float32(size.X), 0)).Add(gtx.Ops)
		}
		stack2.Load()
		op.Offset(f32.Pt(0, float32(size.Y))).Add(gtx.Ops)
	}
	stack1.Load()
	return D{Size: image.Pt(w, h)}
}

func (g *Grid) Event(gtx C) (p *int) {
	columnWidth := g.Width / g.Columns
	rowHeight := g.Height / g.Rows
	for _, e := range gtx.Events(g) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			fmt.Println("click grid")
			x := int(math.Floor(float64(int(e.Position.X) / columnWidth)))
			y := int(math.Floor(float64(int(e.Position.Y) / rowHeight)))
			n := y*g.Columns + x
			p = &n
		}
	}
	return p
}
