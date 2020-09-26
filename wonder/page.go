package main

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/Almanax/wonderwall/wonder/ui"
	"github.com/rs/xid"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"math"
	"time"
)

type Page interface {
	Start(stop <-chan struct{})
	Event(gtx C) interface{}
	Layout(gtx C) D
}

type Transition struct {
	prev, page Page
	reverse    bool
	time       time.Time
}

func (t *Transition) Start(stop <-chan struct{}) {
	t.page.Start(stop)
}

func (t *Transition) Event(gtx layout.Context) interface{} {
	return t.page.Event(gtx)
}

func (t *Transition) Layout(gtx layout.Context) layout.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	prev, page := t.prev, t.page
	if prev != nil {
		if t.reverse {
			prev, page = page, prev
		}
		now := gtx.Now
		if t.time.IsZero() {
			t.time = now
		}
		prev.Layout(gtx)
		cs := gtx.Constraints
		size := layout.FPt(cs.Max)
		max := float32(math.Sqrt(float64(size.X*size.X + size.Y*size.Y)))
		progress := float32(now.Sub(t.time).Seconds()) * 3
		progress = progress * progress // Accelerate
		if progress >= 1 {
			// Stop animation when complete.
			t.prev = nil
		}
		if t.reverse {
			progress = 1 - progress
		}
		diameter := progress * max
		radius := diameter / 2
		op.InvalidateOp{}.Add(gtx.Ops)
		center := size.Mul(.5)
		clipCenter := f32.Point{X: diameter / 2, Y: diameter / 2}
		off := f32.Affine2D{}.Offset(center.Sub(clipCenter))
		op.Affine(off).Add(gtx.Ops)
		clip.RRect{
			Rect: f32.Rectangle{Max: f32.Point{X: diameter, Y: diameter}},
			NE:   radius, NW: radius, SE: radius, SW: radius,
		}.Add(gtx.Ops)
		op.Affine(off.Invert()).Add(gtx.Ops)
		fill{ui.Rgb(0xffffff)}.Layout(gtx)
	}
	return page.Layout(gtx)
}

type pageStack struct {
	pages    []Page
	stopChan chan<- struct{}
}

func (s *pageStack) Len() int {
	return len(s.pages)
}

func (s *pageStack) Current() Page {
	return s.pages[len(s.pages)-1]
}

func (s *pageStack) Pop() {
	s.stop()
	i := len(s.pages) - 1
	prev := s.pages[i]
	s.pages[i] = nil
	s.pages = s.pages[:i]
	if len(s.pages) > 0 {
		s.pages[i-1] = &Transition{
			reverse: true,
			prev:    prev,
			page:    s.Current(),
		}
		s.start()
	}
}

func (s *pageStack) start() {
	stop := make(chan struct{})
	s.stopChan = stop
	s.Current().Start(stop)
}

func (s *pageStack) Swap(p Page) {
	prev := s.pages[len(s.pages)-1]
	s.pages[len(s.pages)-1] = &Transition{
		prev: prev,
		page: p,
	}
	s.start()
}

func (s *pageStack) Push(p Page) {
	if s.stopChan != nil {
		s.stop()
	}
	if len(s.pages) > 0 {
		p = &Transition{
			prev: s.Current(),
			page: p,
		}
	}
	s.pages = append(s.pages, p)
	s.start()
}

func (s *pageStack) stop() {
	close(s.stopChan)
	s.stopChan = nil
}

func (s *pageStack) Clear(p Page) {
	for len(s.pages) > 0 {
		s.Pop()
	}
	s.Push(p)
}

type WallListPage struct {
	env     *Env
	newWall *widget.Clickable
	updates <-chan struct{}

	list *layout.List

	topbar *Topbar
	walls  []interface{}
	clicks []gesture.Click
}

func NewWallListPage(env *Env) *WallListPage {
	clicks := make([]gesture.Click, len(sampleWalls))
	return &WallListPage{
		env:     env,
		newWall: &widget.Clickable{},
		list:    &layout.List{Axis: layout.Vertical},
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
	return nil
}

func (p *WallListPage) Layout(gtx C) D {
	insets := layout.Inset{
		Left:  unit.Dp(16),
		Right: unit.Dp(6),
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return p.topbar.Layout(gtx, layout.Inset{}, ui.Label(theme, unit.Dp(22), "Walls").Layout)
		}),
		layout.Flexed(1, func(gtx C) D {
			return insets.Layout(gtx, func(gtx C) D {
				return p.list.Layout(gtx, len(sampleWalls), p.layoutWallItem)
			})
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

type fill struct {
	color color.RGBA
}

func (f fill) Layout(gtx layout.Context) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: f.color}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: d, Baseline: d.Y}
}
