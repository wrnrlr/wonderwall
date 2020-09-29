package shape

import (
	"gioui.org/f32"
)

type Triangle struct {
}

func (t Triangle) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test
func (t Triangle) Hit(p f32.Point) bool {
	return false
}

func (t Triangle) Offset(p f32.Point) Shape {
	return nil
}

func (t Triangle) Draw(gtx C) {

}
