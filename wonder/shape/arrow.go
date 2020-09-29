package shape

import (
	"gioui.org/f32"
	"image/color"
)

type Arrow struct {
	FillColor   *color.RGBA
	StrokeColor *color.RGBA
	StrokeWidth color.RGBA
}

func (a Arrow) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test
func (a Arrow) Hit(p f32.Point) bool {
	return false
}

func (a Arrow) Offset(p f32.Point) Shape {
	return nil
}

func (a Arrow) Draw(gtx C) {

}
