package main

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/Almanax/wonderwall/wonder/ui"
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
	defer op.Save(gtx.Ops).Load()
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

type fill struct {
	color color.NRGBA
}

func (f fill) Layout(gtx layout.Context) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	paint.ColorOp{Color: f.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: d, Baseline: d.Y}
}

type background struct {
	color color.NRGBA
}

func (b background) Layout(gtx C, w layout.Widget) D {
	dims := w(gtx)
	paint.ColorOp{Color: b.color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}
