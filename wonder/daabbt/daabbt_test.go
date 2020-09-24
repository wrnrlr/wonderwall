package daabbt

import (
	"gioui.org/f32"
	"testing"
)

type rect f32.Rectangle

func (r rect) Bounds() f32.Rectangle {
	return f32.Rectangle(r)
}

func TestDaabbt(t *testing.T) {
	r1 := rect(f32.Rect(10, 10, 20, 20))
	r2 := rect(f32.Rect(10, 10, 50, 50))
	r3 := rect(f32.Rect(60, 60, 80, 100))
	r4 := rect(f32.Rect(60, 60, 100, 80))
	r5 := rect(f32.Rect(300, 300, 600, 400))

	tree := NewTree(f32.Rect(0, 0, 600, 400))
	results := tree.KNearest(f32.Pt(15, 15), 10, func(p f32.Point) bool {
		return true
	})
	assert("expect 0 results", len(results) == 0)

	tree.Insert(r1)
	tree.Insert(r2)
	tree.Insert(r3)
	tree.Insert(r4)
	tree.Insert(r5)

	results = tree.KNearest(f32.Pt(15, 15), 10, func(p f32.Point) bool {
		return true
	})
	assert("expect 2 results", len(results) == 2)
}

func assert(s string, b bool) {
	if !b {
		panic(s)
	}
}
