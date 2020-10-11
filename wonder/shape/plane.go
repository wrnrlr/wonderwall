package shape

import (
	"fmt"
	"gioui.org/f32"
	"github.com/Almanax/wonderwall/wonder/rtree"
	orderedmap "github.com/wk8/go-ordered-map"
)

// A two-dimensional surface that extends infinitely far
type Plane struct {
	Elements *orderedmap.OrderedMap
	Index    rtree.RTree
}

func NewPlane() *Plane {
	return &Plane{
		Elements: orderedmap.New(),
	}
}

func (p *Plane) View(r f32.Rectangle, scale float32, gtx C) {
	fmt.Printf("View: %v %f\n", r, scale)
	p.printElements()
	min := [2]float32{r.Min.X, r.Min.Y}
	max := [2]float32{r.Max.X, r.Max.Y}
	p.Index.Search(min, max, func(min, max [2]float32, value interface{}) bool {
		s, _ := value.(Shape)
		s.Draw(gtx)
		return true
	})
}

func (p *Plane) Insert(s Shape) {
	p.Elements.Set(s.Identity(), s)
	bounds := s.Bounds()
	min, max := [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
	p.Index.Insert(min, max, s)
}

func (p Plane) Within(r f32.Rectangle) Group {
	min, max := [2]float32{r.Min.X, r.Min.Y}, [2]float32{r.Max.X, r.Max.Y}
	p.Index.Search(min, max, func(min [2]float32, max [2]float32, value interface{}) bool {

		return false
	})
	return Group{}
}

func (p Plane) Intersects(r f32.Rectangle) []Shape {
	var results []Shape
	min, max := [2]float32{r.Min.X, r.Min.Y}, [2]float32{r.Max.X, r.Max.Y}
	p.Index.Search(min, max, func(min [2]float32, max [2]float32, value interface{}) bool {
		s, ok := value.(Shape)
		if !ok {
			return false
		}
		results = append(results, s)
		return true
	})
	return results
}

func (p Plane) Hits(pos f32.Point) Shape {
	var result Shape
	min, max := [2]float32{pos.X, pos.Y}, [2]float32{pos.X, pos.Y}
	p.Index.Search(min, max, func(min [2]float32, max [2]float32, value interface{}) bool {
		s, ok := value.(Shape)
		if !ok {
			return false
		}
		if s.Hit(pos) {
			result = s
		}
		return false
	})
	return result
}

func (p *Plane) RemoveAll(ss []Shape) {
	for _, s := range ss {
		p.Elements.Delete(s.Identity())
		bounds := s.Bounds()
		min, max := [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
		p.Index.Delete(min, max, s)
	}
}

func (p Plane) printElements() {
	p.Index.Scan(func(min, max [2]float32, data interface{}) bool {
		s, _ := data.(Shape)
		fmt.Printf("shape: %v\n", s.Bounds())
		return true
	})
}

func intersects(r1, r2 f32.Rectangle) bool {
	if r1.Min.X >= r2.Max.X || r2.Max.X >= r1.Min.X {
		return false
	} else if r1.Min.Y <= r2.Max.X || r2.Max.Y <= r1.Min.X {
		return false
	}
	return true
}
