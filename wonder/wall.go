package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
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

	plane  *shape.Plane
	client Client
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
	macro := op.Record(gtx.Ops)
	d1 := p.toolbar.Layout(gtx)
	toolbar := macro.Stop()
	stack := op.Push(gtx.Ops)
	op.Offset(f32.Pt(0, float32(d1.Size.Y)))
	d2 := p.canvasLayout(gtx)
	stack.Pop()
	toolbar.Add(gtx.Ops)
	return D{Size: image.Pt(d1.Size.X, d1.Size.Y+d2.Size.Y)}
}

func (p *WallPage) canvasLayout(gtx C) D {
	p.canvasEvent(gtx)
	max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	p.plane.View(gtx)
	width := float32(p.toolbar.strokeSize.Value) * p.plane.Scale
	p.pen.Draw(gtx, width, p.toolbar.strokeColor.Color)
	if p.toolbar.Tool == SelectionTool {
		p.selection.Draw(p.plane, gtx)
	}
	return D{Size: max}
}

func (p *WallPage) canvasEvent(gtx C) {
	cons := gtx.Constraints
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: cons.Max}).Add(gtx.Ops)
	for _, e := range gtx.Events(p) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		fmt.Printf("Click on canvas\n")
		if e.Type == pointer.Scroll && (e.Modifiers.Contain(key.ModCommand) || e.Modifiers.Contain(key.ModCtrl)) {
			p.zoom(e.Scroll.Y)
		} else if e.Type == pointer.Scroll {
			p.pan(e.Scroll)
		} else {
			p.toolEvent(e, gtx)
		}
	}
	pointer.InputOp{Tag: p, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release | pointer.Scroll}.Add(gtx.Ops)
}

func (p *WallPage) toolEvent(e pointer.Event, gtx layout.Context) {
	var ev interface{}
	switch p.toolbar.Tool {
	case SelectionTool:
		ev = p.selection.Event(e, p.plane, gtx)
	case PenTool:
		ev = p.pen.Event(e, gtx)
	case TextTool:
		ev = p.text.Event(e, gtx)
	case ImageTool:
		ev = p.images.Event(e, gtx)
	}
	switch ev := ev.(type) {
	case AddLineEvent:
		points := p.transformPoints(ev.Points)
		l := shape.NewPolyline(points, p.toolbar.strokeColor.Color, float32(p.toolbar.strokeSize.Value))
		p.plane.Insert(l)
	case AddTextEvent:
		pos := p.transformPoint(e.Position.Mul(gtx.Metric.PxPerDp))
		txt := shape.NewText(pos.X, pos.Y, "Text", blue, float32(30), theme.Shaper)
		p.plane.Insert(txt)
	case AddImageEvent:
		img := paint.NewImageOp(ev.Image)
		sh := shape.NewImage(e.Position.X, e.Position.Y, &img)
		p.plane.Insert(sh)
	case DragSelectionEvent:
		p.moveSelection(ev.Point)
	case MoveSelectionEvent:
		p.moveSelection(ev.Point)
	}
}

func (p *WallPage) moveSelection(delta f32.Point) {
	for sh, _ := range p.selection.selection {
		sh.Move(delta)
		p.plane.Update(sh)
	}
}

func (p *WallPage) transformPoint(point f32.Point) f32.Point {
	tr := p.plane.GetTransform().Invert()
	return tr.Transform(point)
}

func (p *WallPage) transformPoints(points []f32.Point) []f32.Point {
	tr := p.plane.GetTransform().Invert()
	for i, point := range points {
		points[i] = tr.Transform(point)
	}
	return points
}
