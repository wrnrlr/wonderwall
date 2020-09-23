package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/ui"
	"image"
)

type Toolbar struct {
	Tool      Tool
	theme     *material.Theme
	selection *widget.Clickable
	pen       *widget.Clickable
	text      *widget.Clickable

	strokeSize  *widget.Editor
	strokeColor *widget.Clickable

	delete *widget.Clickable
	undo   *widget.Clickable
	redo   *widget.Clickable
}

func NewToolbar(theme *material.Theme) *Toolbar {
	return &Toolbar{
		Tool:        SelectionTool,
		theme:       theme,
		selection:   new(widget.Clickable),
		pen:         new(widget.Clickable),
		text:        new(widget.Clickable),
		strokeSize:  &widget.Editor{SingleLine: true},
		strokeColor: new(widget.Clickable),
		delete:      new(widget.Clickable),
		undo:        new(widget.Clickable),
		redo:        new(widget.Clickable),
	}
}

func (t *Toolbar) Layout(gtx C) D {
	t.events(gtx)
	stack := layout.Stack{Alignment: layout.NW}
	front := layout.Stacked(func(gtx C) D {
		tools := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
		selection := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.selection, mouseIcon).Layout(gtx)
		})
		pen := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.pen, brushIcon).Layout(gtx)
		})
		text := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.text, textIcon).Layout(gtx)
		})
		strokeSize := layout.Rigid(func(gtx C) D {
			return ui.InputNumber(t.theme, t.strokeSize).Layout(gtx)
		})
		clr := layout.Rigid(func(gtx C) D {
			return ui.Color(t.theme, t.strokeColor).Layout(gtx)
		})
		remove := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.delete, deleteIcon).Layout(gtx)
		})
		undo := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.undo, undoIcon).Layout(gtx)
		})
		redo := layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.redo, redoIcon).Layout(gtx)
		})
		return tools.Layout(gtx, selection, pen, text, strokeSize, clr, remove, undo, redo)
	})
	backg := layout.Expanded(func(gtx C) D {
		cs := gtx.Constraints
		dr := f32.Rectangle{Max: f32.Point{X: float32(cs.Max.X), Y: float32(cs.Min.Y)}}
		paint.ColorOp{Color: t.theme.Color.Primary}.Add(gtx.Ops)
		paint.PaintOp{Rect: dr}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Point{X: cs.Max.X, Y: cs.Min.Y}}
	})
	return stack.Layout(gtx, backg, front)
}

func (t *Toolbar) events(gtx C) {
	if t.selection.Clicked() {
		fmt.Println("clicked selection")
		t.Tool = SelectionTool
	}
	if t.pen.Clicked() {
		fmt.Println("clicked pen")
		t.Tool = PenTool
	}
	if t.text.Clicked() {
		fmt.Println("clicked text")
		t.Tool = TextTool
	}
}
