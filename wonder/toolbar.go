package main

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/Almanax/wonderwall/wonder/colornames"
	"github.com/Almanax/wonderwall/wonder/ui"
	"image"
)

type Toolbar struct {
	Tool      Tool
	theme     *ui.Theme
	selection *widget.Clickable
	pen       *widget.Clickable
	text      *widget.Clickable
	image     *widget.Clickable

	strokeSize  *ui.InputNumberStyle
	strokeColor *ui.ColorPicker

	delete *widget.Clickable
	undo   *widget.Clickable
	redo   *widget.Clickable
	back   *widget.Clickable
}

func NewToolbar(theme *ui.Theme) *Toolbar {
	return &Toolbar{
		Tool:        SelectionTool,
		theme:       theme,
		selection:   new(widget.Clickable),
		pen:         new(widget.Clickable),
		text:        new(widget.Clickable),
		image:       new(widget.Clickable),
		strokeSize:  ui.InputNumber(theme, 10),
		strokeColor: ui.Color(theme, maroon),
		delete:      new(widget.Clickable),
		undo:        new(widget.Clickable),
		redo:        new(widget.Clickable),
		back:        new(widget.Clickable),
	}
}

func (t *Toolbar) Layout(gtx C) D {
	t.events(gtx)

	macro := op.Record(gtx.Ops)
	dims := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.back, backIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.selection, mouseIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.pen, brushIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.text, textIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.image, imageIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return t.strokeSize.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return t.strokeColor.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.delete, deleteIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.undo, undoIcon).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.Item(t.theme, t.redo, redoIcon).Layout(gtx)
		}))
	call := macro.Stop()

	backgroundSize := image.Point{X: gtx.Constraints.Max.X, Y: dims.Size.Y}
	paint.FillShape(gtx.Ops, colornames.Lightgreen, clip.Rect{Min: backgroundSize}.Op())
	call.Add(gtx.Ops)
	return layout.Dimensions{Size: backgroundSize}
}

func (t *Toolbar) events(gtx C) interface{} {
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
	if t.image.Clicked() {
		fmt.Println("clicked text")
		t.Tool = ImageTool
	}
	if t.back.Clicked() {
		fmt.Println("clicked list wall")
		return BackEvent{}
	}
	if t.delete.Clicked() {
		fmt.Println("Delete selection")
		return DeleteEvent{}
	}
	if clr := t.strokeColor.Event(gtx); clr != nil {

	}
	return nil
}
