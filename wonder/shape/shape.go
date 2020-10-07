package shape

import (
	"gioui.org/f32"
)

type Shape interface {
	// Axis aligned bounding box
	Bounds() f32.Rectangle

	// Hit test
	Hit(p f32.Point) bool

	//
	Offset(p f32.Point) Shape

	//
	Draw(gtx C)

	Eq(s2 Shape) bool

	Identity() string
}
