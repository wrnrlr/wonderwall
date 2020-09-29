package shape

import (
	"gioui.org/f32"
)

// A two-dimensional surface that extends infinitely far
type Plane struct {
	Elements Group
}

func (p Plane) View(r f32.Rectangle, gtx C) {
	// Find elements within r
	offset := f32.Pt(r.Dx(), r.Dy())
	for _, s := range p.Elements {
		if intersects(r, s.Bounds()) {
			s.Offset(offset).Draw(gtx)
		}
	}
}

func (p *Plane) Add(s Shape) {
	p.Elements.Append(s)
}

func (p Plane) Within(r f32.Rectangle) Group {
	return Group{}
}

func intersects(r1, r2 f32.Rectangle) bool {
	if r1.Min.X >= r2.Max.X || r2.Max.X >= r1.Min.X {
		return false
	} else if r1.Min.Y <= r2.Max.X || r2.Max.Y <= r1.Min.X {
		return false
	}
	return true
}
