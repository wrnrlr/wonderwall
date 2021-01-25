package main

import "github.com/wrnrlr/wonderwall/wonder/shape"

type Diff struct {
	Prev []shape.Shape
	Next []shape.Shape
}

func (a Diff) Apply(plane shape.Plane) {
	if a.Prev == nil {
		plane.InsertAll(a.Next)
	} else if a.Next == nil {
		plane.RemoveAll(a.Next)
	} else {
		plane.UpdateAll(a.Next)
	}
}

func (a Diff) Invert() Diff {
	return Diff{Prev: a.Next, Next: a.Prev}
}

type State struct {
}
