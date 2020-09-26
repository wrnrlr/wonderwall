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
	exitIcon   = loadIcon(icons.ActionExitToApp)
)

var theme *ui.Theme

func main() {
	theme = ui.MenuTheme(gofont.Collection())
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)), app.Title("Wonderwall"))
		a := NewApp()
		if err := a.loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type Env struct {
	insets layout.Inset
	client *Client
	redraw func()
}

type App struct {
	theme *ui.Theme
	w     *app.Window
	env   Env

	stack pageStack
}

func NewApp() *App {
	theme := ui.MenuTheme(gofont.Collection())
	a := &App{theme: theme}
	a.env.redraw = a.w.Invalidate
	return a
}

func (a *App) loop(w *app.Window) error {
	var updates <-chan struct{}
	var ops op.Ops
	for {
		select {
		case <-updates:
			fmt.Println("Updates...")
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.StageEvent:
				if e.Stage >= system.StageRunning {
					//if a.env.client == nil {
					//	a.env.client = getClient()
					//	updates = a.env.client.register(a)
					//	defer a.env.client.unregister(a)
					//}
					if a.stack.Len() == 0 {
						a.stack.Push(NewLoginPage(&a.env))
					}
				}
			case *system.CommandEvent:
				switch e.Type {
				case system.CommandBack:
					if a.stack.Len() > 1 {
						a.stack.Pop()
						e.Cancel = true
						a.w.Invalidate()
					}
				}
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				a.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func (a *App) Layout(gtx C) {
	a.update(gtx)
	a.stack.Current().Layout(gtx)
}

func (a *App) update(gtx layout.Context) {
	page := a.stack.Current()
	if e := page.Event(gtx); e != nil {
		switch e := e.(type) {
		case BackEvent:
			a.stack.Pop()
		case LoginEvent:
			fmt.Println("LoginEvent")
			//a.env.client.SetAccount(e.Account)
			a.stack.Clear(NewWallListPage(&a.env))
			//a.stack.Clear(NewWallPage(&a.env, xid.New()))
		case ShowWallListEvent:
			fmt.Println("Show Wall List")
			a.stack.Swap(NewWallListPage(&a.env))
		case ShowWallEvent:
			fmt.Println("Show Wall")
			a.stack.Push(NewWallPage(&a.env, e.WallID))
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
