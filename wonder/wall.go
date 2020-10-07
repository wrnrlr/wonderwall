package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
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

	windowSize image.Point
	plane      *shape.Plane
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
		plane:     &shape.Plane{},
	}
}

func (p *WallPage) Start(stop <-chan struct{}) {}

func (p *WallPage) Event(gtx C) interface{} {
	size := gtx.Constraints.Max
	if p.windowSize.X != size.X || p.windowSize.Y != size.Y {
		p.windowSize = size
		//p.tree = daabbt.NewTree(f32.Rect(0, 0, float32(size.X), float32(size.Y)))
	}
	switch p.toolbar.Tool {
	case SelectionTool:
		if e := p.selection.Event(p.plane, gtx); e != nil {
			fmt.Printf("Selection event: %v", e)
		}
	case PenTool:
		if e := p.pen.Event(gtx); e != nil {
			l := &shape.Polyline{Points: e, Width: float32(p.toolbar.strokeSize.Value), Color: p.toolbar.strokeColor.Color}
			p.plane.Insert(l)
			//l.Register(p.tree)
		}
	case TextTool:
		if e := p.text.Event(gtx); e != nil {
			txt := shape.NewText(e.Position.X, e.Position.Y, "Text", blue, float32(30), theme.Shaper)
			p.plane.Insert(txt)
		}
	case ImageTool:
		if e := p.images.Event(gtx); e != nil {
			p.plane.Insert(e)
		}
	default:
	}
	if e := p.toolbar.events(gtx); e != nil {
		return e
	}
	return nil
}

func (p *WallPage) Layout(gtx C) D {
	stack := layout.Stack{}
	toolbar := layout.Stacked(p.toolbar.Layout)
	canvas := layout.Expanded(func(gtx C) D {
		p.Draw(gtx)
		p.pen.Draw(gtx, float32(p.toolbar.strokeSize.Value), p.toolbar.strokeColor.Color)
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		if p.toolbar.Tool == SelectionTool {
			for s, _ := range p.selection.selection {
				b := s.Bounds()
				shape.Rectangle{Rectangle: b, FillColor: nil, StrokeColor: &green, StrokeWidth: unit.Dp(1).V}.Draw(gtx)
			}
		}
		return D{Size: max}
	})
	dims := stack.Layout(gtx, canvas, toolbar)
	//p.tree.Draw(gtx)
	return dims
}

func (p *WallPage) Draw(gtx C) {
	view := f32.Rect(0, 0, 0, 0)
	p.plane.View(view, gtx)
	//for i := range p.lines {
	//	p.lines[i].Draw(gtx)
	//}
	//for i := range p.texts {
	//	p.texts[i].Draw(gtx)
	//}
	//if p.toolbar.Tool == SelectionTool {
	//	for i := range p.lines {
	//		p.lines[i].boxes(gtx)
	//	}
	//}
}
