package shape

import (
	"gioui.org/f32"
	"image/color"
)

type Text struct {
	ID          string
	X, Y        float32
	Text        string
	StrokeColor color.RGBA
	FontWidth   float32
}

func (t Text) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test
func (t Text) Hit(p f32.Point) bool {
	return false
}

func (t Text) Offset(p f32.Point) Shape {
	return nil
}

func (t Text) Draw(gtx C) {
	//defer op.Push(ops).Pop()
	//op.Offset(f32.Point{X: t.X, Y: t.Y - t.FontWidth}).Add(gtx.Ops)
	//l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(t.FontWidth), t.Text)
	//l.Color = t.StrokeColor
	//l.Layout(gtx)
}
