package main

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
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

	plane *shape.Plane
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
		plane:     shape.NewPlane(),
	}
}

func (p *WallPage) Start(stop <-chan struct{}) {}

func (p *WallPage) Event(gtx C) interface{} {
	p.canvasEvent(gtx)
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
	p.plane.Offset = p.plane.Offset.Add(offset)
}

func (p *WallPage) zoom(x float32) {
	const scaleBy = 1.2
	if scaleBy > x {
		p.plane.Scale *= scaleBy
	} else {
		p.plane.Scale /= scaleBy
	}
}

func (p *WallPage) Layout(gtx C) D {
	stack := layout.Stack{}
	toolbar := layout.Stacked(p.toolbar.Layout)
	canvas := layout.Expanded(p.canvasLayout)
	dims := stack.Layout(gtx, canvas, toolbar)
	return dims
}

func (p *WallPage) canvasLayout(gtx C) D {
	p.plane.View(gtx)
	p.pen.Draw(gtx, float32(p.toolbar.strokeSize.Value), p.toolbar.strokeColor.Color)
	max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	if p.toolbar.Tool == SelectionTool {
		p.selection.Draw(p.plane, gtx)
	}
	return D{Size: max}
}

func (p *WallPage) canvasEvent(gtx C) {
	switch p.toolbar.Tool {
	case SelectionTool:
		e := p.selection.Event(p.plane, gtx)
		if e == nil {
			return
		}
		switch e := e.(type) {
		case PanEvent:
			p.pan(e.Offset)
		case ZoomEvent:
			p.zoom(e.Scroll)
		}
	case PenTool:
		e := p.pen.Event(p.plane, gtx)
		if e != nil {
			l := shape.NewPolyline(e, p.toolbar.strokeColor.Color, float32(p.toolbar.strokeSize.Value))
			p.plane.Insert(l)
		}
	case TextTool:
		if e := p.text.Event(p.plane, gtx); e != nil {
			scale := 1 / gtx.Metric.PxPerDp
			pos := e.Position.Mul(scale)
			txt := shape.NewText(pos.X, pos.Y, "Text", blue, float32(30), theme.Shaper)
			p.plane.Insert(txt)
		}
	case ImageTool:
		if e := p.images.Event(p.plane, gtx); e != nil {
			p.plane.Insert(e)
		}
	default:
	}
}

// https://math.stackexchange.com/questions/514212/how-to-scale-a-rectangle
