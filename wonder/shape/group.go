package shape

import (
	"gioui.org/f32"
	"gioui.org/op"
)

type Group []Shape

func (s Group) Offset(p f32.Point) Shape {
	var result Group
	for _, e := range s {
		result = append(result, e.Offset(p))
	}
	return result
}

func (s Group) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test
func (s Group) Hit(p f32.Point) bool {
	return false
}

func (s Group) Draw(ops op.Ops) {

}
