package shape

import (
	"gioui.org/f32"
	"gioui.org/op"
)

type Shape interface {
	// Axis aligned bounding box
	Bounds() f32.Rectangle

	// Hit test
	Hit(p f32.Point) bool

	//
	Offset(p f32.Point) Shape

	//
	Draw(ops op.Ops)
}
