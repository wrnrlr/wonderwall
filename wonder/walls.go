package main

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/Almanax/wonderwall/wonder/ui"
	"github.com/rs/xid"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

type WallListPage struct {
	env     *Env
	newWall *widget.Clickable
	updates <-chan struct{}

	list *layout.List

	topbar  *Topbar
	filters *WallFilter

	add  *widget.Clickable
	user *widget.Clickable

	walls  []interface{}
	clicks []gesture.Click
}

func NewWallListPage(env *Env) *WallListPage {
	clicks := make([]gesture.Click, len(sampleWalls))
	return &WallListPage{
		env:     env,
		newWall: &widget.Clickable{},
		list:    &layout.List{Axis: layout.Vertical},
		topbar:  NewTopbar(false),
		filters: NewWallFilter(),
		add:     &widget.Clickable{},
		user:    &widget.Clickable{},
		clicks:  clicks}
}

func (p *WallListPage) Start(stop <-chan struct{}) {}

func (p *WallListPage) Event(gtx C) interface{} {
	for i := range p.clicks {
		click := &p.clicks[i]
		for _, e := range click.Events(gtx) {
			if e.Type == gesture.TypeClick {
				w := sampleWalls[i]
				return ShowWallEvent{w.ID}
			}
		}
	}
	if p.add.Clicked() {
		return ShowAddWallEvent{}
	}
	if p.user.Clicked() {
		return ShowUserEvent{}
	}
	return nil
}

func (p *WallListPage) Layout(gtx C) D {
	insets := layout.Inset{
		Left:  unit.Dp(16),
		Right: unit.Dp(6),
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return p.topbar.Layout(gtx, layout.Inset{}, p.LayoutMenu)
		}),
		layout.Rigid(func(gtx C) D {
			return p.filters.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return insets.Layout(gtx, func(gtx C) D {
				return p.list.Layout(gtx, len(sampleWalls), p.layoutWallItem)
			})
		}))
}

func (p *WallListPage) LayoutMenu(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return ui.Title(theme, "Walls", gtx)
		}),
		layout.Rigid(ui.Item(theme, p.add, addBoxIcon).Layout),
		layout.Rigid(ui.Item(theme, p.user, userIcon).Layout))
}

func (p *WallListPage) layoutWallItem(gtx C, i int) D {
	click := &p.clicks[i]
	dims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6)}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(ui.H3(theme, sampleWalls[i].Title).Layout),
					layout.Rigid(ui.Caption(theme, "Modified Today").Layout))
			}),
			layout.Rigid(func(gtx C) D {
				ico := (&ui.Icon{Src: icons.NavigationMoreVert, Size: theme.TextSize.Scale(48.0 / 16.0)}).Image(gtx.Metric, theme.Color.Text)
				ico.Add(gtx.Ops)
				paint.PaintOp{Rect: f32.Rectangle{Max: toPointF(ico.Size())}}.Add(gtx.Ops)
				dims := layout.Dimensions{Size: ico.Size()}
				dims.Size.X += gtx.Px(unit.Dp(4))
				return dims
			}))
	})
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	click.Add(gtx.Ops)
	fill{green}.Layout(gtx)
	return dims
}

type Wall struct {
	ID    xid.ID
	Title string
}

var sampleWalls = []*Wall{
	{xid.New(), "Hello, World"},
	{xid.New(), "Wonderwall Sprint"},
	{xid.New(), "Business Model Canvas"},
	{xid.New(), "IT Network"},
	{xid.New(), "Mind Map"},
	{xid.New(), "Inspirational Talk"},
}

type WallFilter struct {
	workspace *widget.Clickable
	text      *widget.Editor
	clear     *gesture.Click
	submit    *gesture.Click
	order     *widget.Clickable
}

func NewWallFilter() *WallFilter {
	return &WallFilter{
		workspace: new(widget.Clickable),
		text:      &widget.Editor{SingleLine: true},
		clear:     new(gesture.Click),
		submit:    new(gesture.Click),
		order:     new(widget.Clickable),
	}
}

func (w *WallFilter) Layout(gtx C) D {
	w.event(gtx)
	stack := layout.Stack{Alignment: layout.NW}
	bg := layout.Expanded(func(gtx C) D {
		cs := gtx.Constraints
		dr := f32.Rectangle{Max: f32.Point{X: float32(cs.Max.X), Y: float32(cs.Min.Y)}}
		paint.ColorOp{Color: theme.Color.Primary}.Add(gtx.Ops)
		paint.PaintOp{Rect: dr}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Point{X: cs.Max.X, Y: cs.Min.Y}}
	})
	fg := layout.Stacked(func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(ui.Item(theme, w.workspace, filterIcon).Layout),
			layout.Flexed(1, ui.InputText(theme, w.text, "Search").Layout),
			layout.Rigid(ui.Item(theme, w.order, sortIcon).Layout))
	})
	return stack.Layout(gtx, bg, fg)
}

func (w *WallFilter) event(gtx C) {
	w.text.Events()
}
