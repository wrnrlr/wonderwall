package shape

import (
	"gioui.org/f32"
	"gioui.org/op"
)

type Group []Shape

func (g Group) Offset(p f32.Point) Shape {
	var result Group
	for _, s := range g {
		result = append(result, s.Offset(p))
	}
	return result
}

func (g Group) Bounds() f32.Rectangle {
	return f32.Rectangle{}
}

// Hit test
func (g Group) Hit(p f32.Point) bool {
	return false
}

func (g Group) Draw(ops op.Ops) {

}

func (g Group) Append(s Shape) Group {
	return append(g, g)
}
