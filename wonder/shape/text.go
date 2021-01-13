package shape

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/rs/xid"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"unicode/utf8"
)

type Text struct {
	ID          string
	X, Y        float32
	Text        string
	StrokeColor color.NRGBA
	FontWidth   float32

	font   text.Font
	shaper text.Shaper
}

func NewText(x, y float32, txt string, col color.NRGBA, width float32, sh text.Shaper) *Text {
	id := xid.New().String()
	return &Text{ID: id, X: x, Y: y, Text: txt, StrokeColor: col, FontWidth: width, shaper: sh}
}

func (t *Text) Bounds() f32.Rectangle {
	textSize := fixed.I(int(t.FontWidth))
	lines := t.shaper.LayoutString(t.font, textSize, inf, t.Text)
	dims := linesDimens(lines)
	//clip := textPadding(lines)
	//clip.Max = clip.Max.Add(dims.Size
	min := f32.Pt(t.X, t.Y-t.FontWidth/2)
	max := min.Add(layout.FPt(dims.Size))
	return f32.Rectangle{Min: min, Max: max}
}

// Hit test
func (t Text) Hit(p f32.Point) bool {
	return true
}

func (t Text) Offset(p f32.Point) Shape {
	return nil
}

func (t Text) Draw(gtx C) {
	defer op.Push(gtx.Ops).Pop()
	scale := gtx.Metric.PxPerDp
	p := f32.Point{X: t.X, Y: t.Y - t.FontWidth/2}.Mul(scale)
	width := t.FontWidth * scale
	op.Offset(p).Add(gtx.Ops)
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(width), t.Text)
	l.Color = t.StrokeColor
	l.Layout(gtx)
}

func (t *Text) Move(delta f32.Point) {
	pos := f32.Point{t.X, t.Y}.Add(delta)
	t.X, t.Y = pos.X, pos.Y
}

// Label stuff

// Label is a widget for laying out and drawing text.
type label struct {
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
}

type lineIterator struct {
	Lines     []text.Line
	Clip      image.Rectangle
	Alignment text.Alignment
	Width     int
	Offset    image.Point

	y, prevDesc fixed.Int26_6
	txtOff      int
}

const inf = 1e6

func (l *lineIterator) Next() (text.Layout, image.Point, bool) {
	for len(l.Lines) > 0 {
		line := l.Lines[0]
		l.Lines = l.Lines[1:]
		x := align(l.Alignment, line.Width, l.Width) + fixed.I(l.Offset.X)
		l.y += l.prevDesc + line.Ascent
		l.prevDesc = line.Descent
		// Align baseline and line start to the pixel grid.
		off := fixed.Point26_6{X: fixed.I(x.Floor()), Y: fixed.I(l.y.Ceil())}
		l.y = off.Y
		off.Y += fixed.I(l.Offset.Y)
		if (off.Y + line.Bounds.Min.Y).Floor() > l.Clip.Max.Y {
			break
		}
		layout := line.Layout
		start := l.txtOff
		l.txtOff += len(line.Layout.Text)
		if (off.Y + line.Bounds.Max.Y).Ceil() < l.Clip.Min.Y {
			continue
		}
		for len(layout.Advances) > 0 {
			_, n := utf8.DecodeRuneInString(layout.Text)
			adv := layout.Advances[0]
			if (off.X + adv + line.Bounds.Max.X - line.Width).Ceil() >= l.Clip.Min.X {
				break
			}
			off.X += adv
			layout.Text = layout.Text[n:]
			layout.Advances = layout.Advances[1:]
			start += n
		}
		end := start
		endx := off.X
		rune := 0
		for n, r := range layout.Text {
			if (endx + line.Bounds.Min.X).Floor() > l.Clip.Max.X {
				layout.Advances = layout.Advances[:rune]
				layout.Text = layout.Text[:n]
				break
			}
			end += utf8.RuneLen(r)
			endx += layout.Advances[rune]
			rune++
		}
		offf := image.Point{X: off.X.Floor(), Y: off.Y.Floor()}
		return layout, offf, true
	}
	return text.Layout{}, image.Point{}, false
}

func (l label) bounds(gtx layout.Context, s text.Shaper, font text.Font, size unit.Value, txt string) layout.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Px(size))
	lines := s.LayoutString(font, textSize, cs.Max.X, txt)
	if max := l.MaxLines; max > 0 && len(lines) > max {
		lines = lines[:max]
	}
	dims := linesDimens(lines)
	dims.Size = cs.Constrain(dims.Size)
	clip := textPadding(lines)
	clip.Max = clip.Max.Add(dims.Size)
	return dims
}

func (l label) Layout(gtx layout.Context, s text.Shaper, font text.Font, size unit.Value, txt string) layout.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Px(size))
	lines := s.LayoutString(font, textSize, cs.Max.X, txt)
	if max := l.MaxLines; max > 0 && len(lines) > max {
		lines = lines[:max]
	}
	dims := linesDimens(lines)
	dims.Size = cs.Constrain(dims.Size)
	cl := textPadding(lines)
	cl.Max = cl.Max.Add(dims.Size)
	it := lineIterator{
		Lines:     lines,
		Clip:      cl,
		Alignment: l.Alignment,
		Width:     dims.Size.X,
	}
	for {
		l, off, ok := it.Next()
		if !ok {
			break
		}
		stack := op.Push(gtx.Ops)
		op.Offset(layout.FPt(off)).Add(gtx.Ops)
		s.Shape(font, textSize, l).Add(gtx.Ops)
		clip.Rect(cl.Sub(off)).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		stack.Pop()
	}
	return dims
}

func textPadding(lines []text.Line) (padding image.Rectangle) {
	if len(lines) == 0 {
		return
	}
	first := lines[0]
	if d := first.Ascent + first.Bounds.Min.Y; d < 0 {
		padding.Min.Y = d.Ceil()
	}
	last := lines[len(lines)-1]
	if d := last.Bounds.Max.Y - last.Descent; d > 0 {
		padding.Max.Y = d.Ceil()
	}
	if d := first.Bounds.Min.X; d < 0 {
		padding.Min.X = d.Ceil()
	}
	if d := first.Bounds.Max.X - first.Width; d > 0 {
		padding.Max.X = d.Ceil()
	}
	return
}

func linesDimens(lines []text.Line) layout.Dimensions {
	var width fixed.Int26_6
	var h int
	var baseline int
	if len(lines) > 0 {
		baseline = lines[0].Ascent.Ceil()
		var prevDesc fixed.Int26_6
		for _, l := range lines {
			h += (prevDesc + l.Ascent).Ceil()
			prevDesc = l.Descent
			if l.Width > width {
				width = l.Width
			}
		}
		h += lines[len(lines)-1].Descent.Ceil()
	}
	w := width.Ceil()
	return layout.Dimensions{
		Size: image.Point{
			X: w,
			Y: h,
		},
		Baseline: h - baseline,
	}
}

func align(align text.Alignment, width fixed.Int26_6, maxWidth int) fixed.Int26_6 {
	mw := fixed.I(maxWidth)
	switch align {
	case text.Middle:
		return fixed.I(((mw - width) / 2).Floor())
	case text.End:
		return fixed.I((mw - width).Floor())
	case text.Start:
		return 0
	default:
		panic(fmt.Errorf("unknown alignment %v", align))
	}
}

func (t *Text) Eq(s Shape) bool {
	t2, ok := s.(*Text)
	if !ok {
		return false
	}
	return t.ID == t2.ID
}

func (t *Text) Identity() string {
	return t.ID
}
