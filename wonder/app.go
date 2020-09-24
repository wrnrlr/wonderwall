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
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Almanax/wonderwall/wonder/daabbt"
	"github.com/Almanax/wonderwall/wonder/ui"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"log"
	"os"
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

type App struct {
	theme         *material.Theme
	disabledTheme *material.Theme
	activeTheme   *material.Theme

	toolbar *Toolbar

	selection *Selection
	pen       *Pen
	text      *Text

	penConfig *PenConfig

	windowSize image.Point
	tree       *daabbt.Node
	lines      []*Line
	texts      []Text
}

func NewApp() *App {
	theme := ui.CustomTheme(gofont.Collection())
	penConfig := &PenConfig{StrokeSize: 10, StrokeColor: maroon}
	return &App{
		theme:     theme,
		toolbar:   NewToolbar(theme),
		selection: new(Selection),
		pen:       new(Pen),
		text:      new(Text),
		penConfig: penConfig,
		tree:      nil,
		lines:     nil,
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
	size := gtx.Constraints.Max
	if a.windowSize.X != size.X || a.windowSize.Y != size.Y {
		a.windowSize = size
		a.tree = daabbt.NewTree(f32.Rect(0, 0, float32(size.X), float32(size.Y)))
	}
	switch a.toolbar.Tool {
	case SelectionTool:
		if e := a.selection.Event(a.tree, gtx); e != nil {
			fmt.Printf("Selection event: %v", e)
		}
	case PenTool:
		if e := a.pen.Event(gtx); e != nil {
			l := &Line{Points: e, Width: float32(a.penConfig.StrokeSize)}
			a.lines = append(a.lines, l)
			l.Register(a.tree)
		}
	case TextTool:
	default:
	}
}

func (a *App) Layout(gtx C) D {
	stack := layout.Stack{}
	toolbar := layout.Stacked(a.toolbar.Layout)
	canvas := layout.Expanded(func(gtx C) D {
		a.Draw(gtx)
		a.pen.Draw(gtx, a.penConfig)
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		return D{Size: max}
	})
	dims := stack.Layout(gtx, canvas, toolbar)
	a.tree.Draw(gtx)
	return dims
}

func (a *App) Draw(gtx C) {
	for i := range a.lines {
		a.lines[i].Draw(a.penConfig, gtx)
	}
	if a.toolbar.Tool == SelectionTool {
		for i := range a.lines {
			a.lines[i].boxes(gtx)
		}
	}
}

type Selection struct {
	events []f32.Point
}

func (s *Selection) Event(tree *daabbt.Node, gtx C) []f32.Point {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(s) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			results := tree.KNearest(e.Position, 10, func(p f32.Point) bool {
				return true
			})
			fmt.Printf("results: %v\n", results)
		case pointer.Drag:
		case pointer.Release, pointer.Cancel:
		}
	}
	pointer.InputOp{Tag: s, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	return nil
}

type Text struct {
	pointer gesture.Click
}

type Canvas struct{}

func (c Canvas) Layout(gtx C) D {
	return D{Size: gtx.Constraints.Min}
}

func (c Canvas) Events(gtx C) {

}
