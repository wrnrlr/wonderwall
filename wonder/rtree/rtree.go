package rtree

import (
	"gioui.org/f32"
	"image"
)

type Shape interface {
	Bounds() image.Rectangle
	Area() float32
}

type rect struct {
	f32.Rectangle
	data interface{}
}

type RTree struct {
	height   int
	root     rect
	count    int
	reinsert []rect
}

func (t *RTree) Insert(s Shape) {

}

func (t *RTree) Remove(s Shape) {}

func (t *RTree) Within(r f32.Rectangle) {}

func (t *RTree) Near(p f32.Point) {}
