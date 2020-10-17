package shape

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/op"
	"github.com/Almanax/wonderwall/wonder/rtree"
	orderedmap "github.com/wk8/go-ordered-map"
	"golang.org/x/image/colornames"
)

// A two-dimensional surface that extends infinitely far
type Plane struct {
	Elements *orderedmap.OrderedMap
	Index    rtree.RTree
	Offset   f32.Point
	Scale    float32
	Width    float32
	Height   float32
}

func NewPlane() *Plane {
	return &Plane{
		Elements: orderedmap.New(),
		Scale:    1,
	}
}

func (p *Plane) View(gtx C) {
	//fmt.Printf("View: %v %f\n", r, scale)
	//p.printElements()
	cons := gtx.Constraints
	p.Width, p.Height = float32(cons.Max.X), float32(cons.Max.Y)

	center := p.Center()
	scaledWidth, scaledHeight := p.Width*p.Scale, p.Height*p.Scale
	min := [2]float32{center.X - scaledWidth/2, center.Y - scaledHeight/2}
	max := [2]float32{center.X + scaledWidth/2, center.Y + scaledHeight/2}

	tr := p.GetTransform()
	defer op.Push(gtx.Ops).Pop()
	op.Affine(tr).Add(gtx.Ops)

	p.Index.Search(min, max, func(min, max [2]float32, key interface{}) bool {
		value, _ := p.Elements.Get(key)
		s, _ := value.(Shape)
		s.Draw(gtx)
		return true
	})

	for pair := p.Elements.Oldest(); pair != nil; pair = pair.Next() {
		s, _ := pair.Value.(Shape)
		Rectangle{s.Bounds(), nil, &colornames.Lightgray, float32(1)}.Draw(gtx)
	}
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
	p.Index.Search(min, max, func(min [2]float32, max [2]float32, key interface{}) bool {
		value, found := p.Elements.Get(key)
		if !found {
			return false
		}
		s, ok := value.(Shape)
		if !ok {
			return false
		}
		if s.Hit(pos) {
			result = s
			return true
		}
		return false
	})
	return result
}

func (p *Plane) Insert(s Shape) {
	p.Elements.Set(s.Identity(), s)
	bounds := s.Bounds()
	min, max := [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
	p.Index.Insert(min, max, s.Identity())
}

func (p *Plane) InsertAll(ss []Shape) {
	for _, s := range ss {
		p.Insert(s)
	}
}

func (p *Plane) Update(s Shape) {
	old, found := p.Elements.Get(s.Identity())
	if !found {
		return
	}
	olds := old.(Shape)
	bounds := olds.Bounds()
	min, max := [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
	p.Index.Delete(min, max, s.Identity())

	p.Elements.Set(s.Identity(), s)
	bounds = s.Bounds()
	min, max = [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
	p.Index.Insert(min, max, s.Identity())
}

func (p *Plane) UpdateAll(ss []Shape) {
	for _, s := range ss {
		p.Update(s)
	}
}

func (p *Plane) Remove(s Shape) {
	p.Elements.Delete(s.Identity())
	bounds := s.Bounds()
	min, max := [2]float32{bounds.Min.X, bounds.Min.Y}, [2]float32{bounds.Max.X, bounds.Max.Y}
	p.Index.Delete(min, max, s.Identity())
}

func (p *Plane) RemoveAll(ss []Shape) {
	for _, s := range ss {
		p.Remove(s)
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

func (p Plane) RelativePoint(point f32.Point, gtx C) f32.Point {
	return point
}

func (p Plane) Center() f32.Point {
	return f32.Pt(p.Offset.X+p.Width/2, p.Offset.Y+p.Height/2)
}

func (p Plane) GetTransform() f32.Affine2D {
	return f32.Affine2D{}.Scale(p.Center(), f32.Pt(p.Scale, p.Scale))
}
