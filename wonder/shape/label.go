package shape

import (
	"gioui.org/f32"
)

type Label struct{}

func (l Label) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

func (l Label) Offset(p f32.Point) Shape {
	return nil
}

// Hit test
func (l Label) Hit(p f32.Point) bool {
	return false
}

func (l Label) Draw(gtx C) {

}
