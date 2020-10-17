package shape

import (
	"gioui.org/f32"
)

type Group struct {
	ID       string
	Elements []Shape
}

func (g *Group) Offset(p f32.Point) Shape {
	var result Group
	for _, s := range g.Elements {
		result.Elements = append(result.Elements, s.Offset(p))
	}
	return &result
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

func (g *Group) Move(delta f32.Point) {}

func (g *Group) Append(s Shape) {
	g.Elements = append(g.Elements, s)
}

func (g *Group) Eq(s2 Shape) bool {
	g2, ok := s2.(*Group)
	if !ok {
		return false
	}
	return g.ID == g2.ID
}

func (g *Group) Identity() string {
	return g.ID
}
