package main

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/shape"
	"github.com/Almanax/wonderwall/wonder/ui"
	"github.com/rs/xid"
	"image"
)

type WallPage struct {
	env    *Env
	WallID xid.ID

	disabledTheme *material.Theme
	activeTheme   *material.Theme

	toolbar *Toolbar

	selection *Selection
	pen       *Pen
	text      *TextWriter

	images ImageService

	offset f32.Point
	scale  float32
	plane  *shape.Plane
}

func NewWallPage(env *Env, wallID xid.ID) *WallPage {
	theme := ui.MenuTheme(gofont.Collection())
	return &WallPage{
		env:       env,
		WallID:    wallID,
		toolbar:   NewToolbar(theme),
		selection: NewSelection(),
		pen:       new(Pen),
		text:      new(TextWriter),
		offset:    f32.Pt(0, 0),
		scale:     1,
		plane:     shape.NewPlane(),
	}
}

func (p *WallPage) Start(stop <-chan struct{}) {}

func (p *WallPage) Event(gtx C) interface{} {
	//size := gtx.Constraints.Max
	//if p.windowSize.X != size.X || p.windowSize.Y != size.Y {
	//	p.windowSize = size
	//	//p.tree = daabbt.NewTree(f32.Rect(0, 0, float32(size.X), float32(size.Y)))
	//}
	switch p.toolbar.Tool {
	case SelectionTool:
		e := p.selection.Event(p.plane, gtx)
		if e == nil {
			return nil
		}
		switch e := e.(type) {
		case PanEvent:
			p.pan(e.Offset)
		case ZoomEvent:
			p.zoom(e.Scroll)
		}
	case PenTool:
		if e := p.pen.Event(gtx); e != nil {
			l := shape.NewPolyline(e, p.toolbar.strokeColor.Color, float32(p.toolbar.strokeSize.Value))
			p.plane.Insert(l)
			//l.Register(p.tree)
		}
	case TextTool:
		if e := p.text.Event(gtx); e != nil {
			scale := 1 / gtx.Metric.PxPerDp
			pos := e.Position.Mul(scale)
			txt := shape.NewText(pos.X, pos.Y, "Text", blue, float32(30), theme.Shaper)
			p.plane.Insert(txt)
		}
	case ImageTool:
		if e := p.images.Event(gtx); e != nil {
			p.plane.Insert(e)
		}
	default:
	}
	if e := p.toolbar.events(gtx); e != nil {
		switch e.(type) {
		case DeleteEvent:
			p.DeleteSelection()
		default:
			return e
		}
	}
	return nil
}

func (p *WallPage) DeleteSelection() {
	p.plane.RemoveAll(p.selection.Elements())
	p.selection.Clear()
}

func (p *WallPage) pan(offset f32.Point) {
	p.offset = p.offset.Add(offset)
}

func (p *WallPage) zoom(x float32) {
	const scaleBy = 1.1
	if scaleBy > x {
		p.scale *= scaleBy
	} else {
		p.scale /= scaleBy
	}
}

func (p *WallPage) Layout(gtx C) D {
	stack := layout.Stack{}
	toolbar := layout.Stacked(p.toolbar.Layout)
	canvas := layout.Expanded(func(gtx C) D {

		cons := gtx.Constraints
		r := f32.Rectangle{Min: f32.Point{X: 0, Y: 0}, Max: layout.FPt(cons.Max)}
		scale := p.scale
		width := r.Dx()
		height := r.Dy()
		scaledWidth := width * scale
		scaledHeight := height * scale
		centerX := r.Min.X + width/2
		centerY := r.Min.Y + height/2
		center := f32.Pt(centerX, centerY)
		minX := r.Min.X + width/2 - scaledWidth/2
		minY := r.Min.Y + height/2 - scaledHeight/2
		maxX := r.Max.X + width/2 - scaledWidth/2
		maxY := r.Max.Y + height/2 - scaledHeight/2
		op.InvalidateOp{}.Add(gtx.Ops)
		defer op.Push(gtx.Ops).Pop()
		tr := f32.Affine2D{}
		tr = tr.Scale(center, f32.Pt(scale, scale))
		op.Affine(tr).Add(gtx.Ops)
		view := f32.Rect(minX, minY, maxX, maxY)
		p.plane.View(view, scale, gtx)
		p.pen.Draw(gtx, float32(p.toolbar.strokeSize.Value), p.toolbar.strokeColor.Color)
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		if p.toolbar.Tool == SelectionTool {
			for pair := p.plane.Elements.Oldest(); pair != nil; pair = pair.Next() {
				s, _ := pair.Value.(shape.Shape)
				b := s.Bounds()
				shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &lightpink, StrokeWidth: unit.Dp(1).V}.Draw(gtx)
			}
			p.plane.Index.Scan(func(min, max [2]float32, data interface{}) bool {
				b := f32.Rect(min[0], min[1], max[0], max[1])
				shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &lightgrey, StrokeWidth: unit.Dp(1).V}.Draw(gtx)
				return true
			})
			for s, _ := range p.selection.selection {
				b := s.Bounds()
				shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &green, StrokeWidth: unit.Dp(1).V}.Draw(gtx)
			}
		}
		return D{Size: max}
	})
	dims := stack.Layout(gtx, canvas, toolbar)
	return dims
}

// https://math.stackexchange.com/questions/514212/how-to-scale-a-rectangle
//func (p *WallPage) Draw(gtx C) {
//	cons := gtx.Constraints
//	scale := p.scale
//	min := f32.Point{X: 0, Y: 0} //.Add(p.offset)
//	max := layout.FPt(cons.Max)  //.Add(p.offset)
//	view := f32.Rectangle{Min: min, Max: max}
//	//for i := range p.lines {
//	//	p.lines[i].Draw(gtx)
//	//}
//	//for i := range p.texts {
//	//	p.texts[i].Draw(gtx)
//	//}
//	//if p.toolbar.Tool == SelectionTool {
//	//	for i := range p.lines {
//	//		p.lines[i].boxes(gtx)
//	//	}
//	//}
//}
