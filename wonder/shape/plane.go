package shape

import (
	"gioui.org/f32"
)

// A two-dimensional surface that extends infinitely far
type Plane struct {
	Elements []Shape
}

func (p Plane) View(r f32.Rectangle, gtx C) {
	// Find elements within r

	offset := f32.Pt(r.Dx(), r.Dy())

	// Translate shape's position in the plane to screen's position
	for _, s := range p.Elements {
		s.Offset(offset)
	}
}

func (p Plane) Within(r f32.Rectangle) Group {

}
