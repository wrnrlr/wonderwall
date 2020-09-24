// Based of https://github.com/asim/quadtree/blob/master/quadtree.go
package daabbt

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image/color"
)

var (
	Capacity = 8
	MaxDepth = 6
)

type Filter func(f32.Point) bool

type Bounded interface {
	Bounds() f32.Rectangle
}

func rectangleContains(b f32.Rectangle, p f32.Point) bool {
	return !(p.X < b.Min.X || p.Y < b.Min.Y || p.X > b.Max.X || p.Y > b.Max.Y)
}

func rectangleOverlaps(b1, b2 f32.Rectangle) bool {
	return !(b1.Max.X < b2.Min.X || b1.Max.Y < b2.Min.Y || b1.Min.X > b2.Max.X || b1.Min.Y > b2.Max.Y)
}

type Node struct {
	boundary f32.Rectangle
	items    []Bounded
	parent   *Node
	children [4]*Node
	depth    int
}

func NewTree(dims f32.Rectangle) *Node {
	return &Node{boundary: dims}
}

func (n *Node) Insert(v Bounded) bool {
	b := v.Bounds()
	if !rectangleOverlaps(n.boundary, b) {
		return false
	}
	if n.isEmpty() {
		if len(n.items) < Capacity {
			n.items = append(n.items, v)
			return true
		}
		if n.depth < MaxDepth {
			n.quarter()
		} else {
			n.items = append(n.items, v)
			return true
		}
	}
	for _, node := range n.children {
		if node.Insert(v) {
			return true
		}
	}
	return false
}

func (n *Node) isEmpty() bool {
	return n.children[0] == nil
}

func (n *Node) Remove(b Bounded) bool {
	if !rectangleOverlaps(n.boundary, b.Bounds()) {
		return false
	}
	if n.children[0] == nil {
		for i, item := range n.items {
			if item != b {
				continue
			}
			if last := len(n.items) - 1; i == last {
				n.items = n.items[:last]
			} else {
				n.items[i] = n.items[last]
				n.items = n.items[:last]
			}
			return true
		}
		return false
	}
	for _, item := range n.children {
		if item.Remove(b) {
			return true
		}
	}
	return false
}

// Split into four
func (n *Node) quarter() {
	if !n.isEmpty() {
		return
	}
	min, max := n.boundary.Min, n.boundary.Max
	halfX, halfY := (max.X-min.X)/2, (max.Y-min.Y)/2
	centerX, centerY := max.X-halfX, max.Y-halfY
	nw := f32.Rect(min.X, min.Y, centerX, centerY)
	n.children[0] = &Node{boundary: nw, depth: n.depth, parent: n}
	ne := f32.Rect(centerX, min.Y, max.X, centerY)
	n.children[1] = &Node{boundary: ne, depth: n.depth, parent: n}
	se := f32.Rect(centerX, centerY, max.X, max.Y)
	n.children[2] = &Node{boundary: se, depth: n.depth, parent: n}
	sw := f32.Rect(min.X, centerY, centerX, max.Y)
	n.children[3] = &Node{boundary: sw, depth: n.depth, parent: n}
	for _, p := range n.items {
		for _, c := range n.children {
			if c.Insert(p) {
				break
			}
		}
	}
	n.items = nil
}

func (n *Node) KNearest(point f32.Point, k int, fn Filter) []Bounded {
	v := make(map[*Node]bool)
	return n.kNearestRoot(point, k, v, fn)
}

func (n *Node) kNearestRoot(p f32.Point, k int, v map[*Node]bool, fn Filter) (results []Bounded) {
	if !rectangleContains(n.boundary, p) {
		return results
	}
	// hit the leaf
	if n.children[0] == nil {
		results = append(results, n.knearest(p, k, v, fn)...)
		if len(results) >= k {
			results = results[:k]
		}
		return results
	}
	for _, node := range n.children {
		results = append(results, node.kNearestRoot(p, k, v, fn)...)
		if len(results) >= k {
			return results[:k]
		}
	}
	if len(results) >= k {
		results = results[:k]
	}
	return results
}

func (n *Node) knearest(p f32.Point, k int, v map[*Node]bool, fn Filter) (results []Bounded) {
	if _, ok := v[n]; ok {
		return results
	} else {
		v[n] = true
	}
	if !rectangleContains(n.boundary, p) {
		return results
	}
	for _, item := range n.items {
		if rectangleContains(item.Bounds(), p) {
			results = append(results, item)
		}
		if len(results) >= k {
			return results[:k]
		}
	}
	if n.children[0] != nil {
		for _, node := range n.children {
			results = append(results, node.knearest(p, k, v, fn)...)
			if len(results) >= k {
				return results[:k]
			}
		}
		if len(results) >= k {
			results = results[:k]
		}
	}
	if n.parent == nil {
		return results
	}
	results = append(results, n.parent.knearest(p, k, v, fn)...)
	if len(results) >= k {
		results = results[:k]
	}
	return results
}

func (n *Node) Draw(gtx layout.Context) {
	shape.Rectangle(n.boundary).Stroke(color.RGBA{255, 182, 193, 255}, float32(2), gtx)
	if n.children[0] != nil {
		n.children[0].Draw(gtx)
		n.children[1].Draw(gtx)
		n.children[2].Draw(gtx)
		n.children[3].Draw(gtx)
	}
}
