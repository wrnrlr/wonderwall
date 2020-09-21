package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/ui"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/wrnrlr/polyline"
)

var (
	maroon    = color.RGBA{127, 0, 0, 255}
	lightgrey = color.RGBA{100, 100, 100, 255}
	black     = color.RGBA{0, 0, 0, 255}
	red       = color.RGBA{255, 0, 0, 255}
	green     = color.RGBA{0, 255, 0, 255}
	blue      = color.RGBA{0, 0, 255, 255}
)

type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	mouseIcon  = loadIcon(icons.ContentSelectAll)
	brushIcon  = loadIcon(icons.ImageBrush)
	textIcon   = loadIcon(icons.EditorTitle)
	deleteIcon = loadIcon(icons.ActionDelete)
	undoIcon   = loadIcon(icons.ContentUndo)
	redoIcon   = loadIcon(icons.ContentRedo)
)

func loadIcon(b []byte) *widget.Icon {
	icon, err := widget.NewIcon(b)
	if err != nil {
		panic(err)
	}
	return icon
}

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		a := NewApp()
		if err := a.loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type Lines struct {
	lines []Line
}

type Line struct {
	points []f32.Point
}

type App struct {
	theme         *material.Theme
	disabledTheme *material.Theme
	activeTheme   *material.Theme

	toolbar *Toolbar

	tool Tool

	pen  *Pen
	text *Text

	penConfig *PenConfig

	lines []Line
	texts []Text
}

func NewApp() *App {
	theme := ui.CustomTheme(gofont.Collection())
	penConfig := &PenConfig{StrokeSize: 10, StrokeColor: maroon}
	return &App{
		theme:     theme,
		toolbar:   NewToolbar(theme),
		tool:      PenTool,
		pen:       new(Pen),
		text:      new(Text),
		penConfig: penConfig,
	}
}

func (a *App) loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			a.Event(gtx)
			a.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) Event(gtx C) {
	switch a.toolbar.Tool {
	case SelectionTool:
	case PenTool:
		fmt.Println("Check pen events")
		if e := a.pen.Event(gtx); e != nil {
			fmt.Println("new line")
			a.lines = append(a.lines, Line{e})
		}
	case TextTool:
	default:
	}
}

func (a *App) Layout(gtx C) D {
	stack := layout.Stack{}
	toolbar := layout.Stacked(a.toolbar.Layout)
	canvas := layout.Expanded(func(gtx C) D {
		for _, l := range a.lines {
			polyline.Draw(l.points, unit.Dp(float32(a.penConfig.StrokeSize)).V, a.penConfig.StrokeColor, gtx)
		}
		a.pen.Draw(gtx, a.penConfig)
		return D{Size: gtx.Constraints.Max}
	})
	return stack.Layout(gtx, canvas, toolbar)
}

type Tool int

const (
	NoTool Tool = iota
	SelectionTool
	PenTool
	TextTool
)

func (t Tool) String() string {
	switch t {
	case NoTool:
		return "NoTool"
	case SelectionTool:
		return "SelectionTool"
	case PenTool:
		return "PenTool"
	case TextTool:
		return "TextTool"
	default:
		return "UnknownTool"
	}
}

type Pen struct {
	events []f32.Point
	grab   bool
}

func (p *Pen) Draw(gtx C, conf *PenConfig) {
	if p.events != nil {
		polyline.Draw(p.events, float32(conf.StrokeSize), maroon, gtx)
	}
}

func (p *Pen) Event(gtx C) []f32.Point {
	var l []f32.Point
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(p) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			p.events = append(p.events, e.Position)
		case pointer.Drag:
			p.events = append(p.events, e.Position)
		case pointer.Release, pointer.Cancel:
			l = append(p.events, e.Position)
			p.events = nil
		}
	}
	p.Add(gtx.Ops)
	return l
}

func (p *Pen) Add(ops *op.Ops) {
	pointer.InputOp{
		Tag:   p,
		Grab:  p.grab,
		Types: pointer.Press | pointer.Drag | pointer.Release,
	}.Add(ops)
}

type PenConfig struct {
	StrokeSize  int
	StrokeColor color.RGBA
}

type Text struct {
	pointer gesture.Click
}

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

type Canvas struct {
}

func (c Canvas) Layout(gtx C) D {
	return D{Size: gtx.Constraints.Min}
}

func (c Canvas) Events(gtx C) {

}

type Minimap struct {
}

func (m Minimap) Layout(gtx C) D {
	return D{Size: gtx.Constraints.Min}
}

func (m Minimap) Events(gtx C) {

}
