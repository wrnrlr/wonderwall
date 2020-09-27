package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/daabbt"
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

	penConfig *PenConfig

	windowSize image.Point
	tree       *daabbt.Node
	lines      []*Line
	texts      []*Text
}

func NewWallPage(env *Env, wallID xid.ID) *WallPage {
	theme := ui.MenuTheme(gofont.Collection())
	penConfig := &PenConfig{StrokeSize: 10, StrokeColor: maroon}
	return &WallPage{
		env:       env,
		WallID:    wallID,
		toolbar:   NewToolbar(theme),
		selection: new(Selection),
		pen:       new(Pen),
		text:      new(TextWriter),
		penConfig: penConfig,
		tree:      nil,
		lines:     nil,
	}
}

func (p *WallPage) Start(stop <-chan struct{}) {}

func (p *WallPage) Event(gtx C) interface{} {
	size := gtx.Constraints.Max
	if p.windowSize.X != size.X || p.windowSize.Y != size.Y {
		p.windowSize = size
		p.tree = daabbt.NewTree(f32.Rect(0, 0, float32(size.X), float32(size.Y)))
	}
	switch p.toolbar.Tool {
	case SelectionTool:
		if e := p.selection.Event(p.tree, gtx); e != nil {
			fmt.Printf("Selection event: %v", e)
		}
	case PenTool:
		if e := p.pen.Event(gtx); e != nil {
			l := &Line{Points: e, StrokeWidth: float32(p.penConfig.StrokeSize)}
			p.lines = append(p.lines, l)
			l.Register(p.tree)
		}
	case TextTool:
		if e := p.text.Event(gtx); e != nil {
			p.texts = append(p.texts, &Text{e.Position, "Text", blue, float32(50)})
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
		p.pen.Draw(gtx, p.penConfig)
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		return D{Size: max}
	})
	dims := stack.Layout(gtx, canvas, toolbar)
	p.tree.Draw(gtx)
	return dims
}

func (p *WallPage) Draw(gtx C) {
	for i := range p.lines {
		p.lines[i].Draw(p.penConfig, gtx)
	}
	for i := range p.texts {
		p.texts[i].Draw(gtx)
	}
	if p.toolbar.Tool == SelectionTool {
		for i := range p.lines {
			p.lines[i].boxes(gtx)
		}
	}
}
