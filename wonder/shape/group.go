package shape

import (
	"gioui.org/f32"
)

type Group struct {
	Elements []Shape
}

func (g Group) Offset(p f32.Point) Shape {
	var result Group
	for _, s := range g.Elements {
		result.Elements = append(result.Elements, s.Offset(p))
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

func (g Group) Draw(gtx C) {
	for _, s := range g.Elements {
		s.Draw(gtx)
	}
}

func (g *Group) Append(s Shape) {
	g.Elements = append(g.Elements, s)
}
