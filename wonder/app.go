package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

type App struct {
	theme         *material.Theme
	disabledTheme *material.Theme
	activeTheme   *material.Theme

	toolbar *Toolbar

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
		pen:       new(Pen),
		text:      new(Text),
		penConfig: penConfig,
		lines:     []Line{{[]f32.Point{f32.Pt(150, 150), f32.Pt(500, 500)}}},
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
		//for _, l := range a.lines {
		//	shape.Line{Points: l.points, Color: a.penConfig.StrokeColor, Width:float32(a.penConfig.StrokeSize)}.Hit(gtx)
		//}
	case PenTool:
		if e := a.pen.Event(gtx); e != nil {
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
		a.Draw(gtx)
		a.pen.Draw(gtx, a.penConfig)
		max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		return D{Size: max}
	})
	dims := stack.Layout(gtx, canvas, toolbar)
	if a.toolbar.Tool == SelectionTool {
		for i := range a.lines {
			a.lines[i].Event(a.penConfig, gtx)
		}
	}
	return dims
}

func (a *App) Draw(gtx C) {
	for i := range a.lines {
		a.lines[i].Draw(a.penConfig, gtx)
	}
}

type Selection struct{}

type Text struct {
	pointer gesture.Click
}

type Canvas struct{}

func (c Canvas) Layout(gtx C) D {
	return D{Size: gtx.Constraints.Min}
}

func (c Canvas) Events(gtx C) {

}
