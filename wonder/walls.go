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

	add  gesture.Click
	user gesture.Click

	walls  []interface{}
	clicks []gesture.Click
}

func NewWallListPage(env *Env) *WallListPage {
	clicks := make([]gesture.Click, len(sampleWalls))
	return &WallListPage{
		env:     env,
		newWall: &widget.Clickable{},
		list:    &layout.List{Axis: layout.Vertical},
		filters: NewWallFilter(),
		topbar:  &Topbar{},
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
	for _, e := range p.add.Events(gtx) {
		if e.Type == gesture.TypeClick {
			return ShowAddWallEvent{}
		}
	}
	for _, e := range p.user.Events(gtx) {
		if e.Type == gesture.TypeClick {
			return ShowUserEvent{}
		}
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
			return ui.Label(theme, unit.Dp(22), "Walls").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			ico := (&ui.Icon{Src: icons.ContentAddBox, Size: unit.Dp(24)}).Image(gtx.Metric, theme.Color.Text)
			ico.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: toPointF(ico.Size())}}.Add(gtx.Ops)
			dims := layout.Dimensions{Size: ico.Size()}
			dims.Size.X += gtx.Px(unit.Dp(4))
			pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
			p.add.Add(gtx.Ops)
			return dims
		}),
		layout.Rigid(func(gtx C) D {
			ico := (&ui.Icon{Src: icons.SocialPerson, Size: unit.Dp(24)}).Image(gtx.Metric, theme.Color.Text)
			ico.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: toPointF(ico.Size())}}.Add(gtx.Ops)
			dims := layout.Dimensions{Size: ico.Size()}
			dims.Size.X += gtx.Px(unit.Dp(4))
			pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
			p.user.Add(gtx.Ops)
			return dims
		}))
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
				//pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
				//click.Add(gtx.Ops)
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
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(ui.Item(theme, w.workspace, filterIcon).Layout),
		layout.Flexed(1, ui.InputText(theme, w.text, "Search").Layout),
		layout.Rigid(ui.Item(theme, w.order, sortIcon).Layout))
}

func (w *WallFilter) event(gtx C) {
	w.text.Events()
}
