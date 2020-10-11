package rtree

import (
	"fmt"
	"testing"
)

func TestRTree(t *testing.T) {
	tree := RTree{}

	tree.Insert([2]float32{-10, -10}, [2]float32{100, 100}, "A")
	tree.Insert([2]float32{50, 50}, [2]float32{70, 120}, "B")
	tree.Insert([2]float32{130, 130}, [2]float32{140, 140}, "C")

	testResultLength(t, tree, [2]float32{0, 0}, [2]float32{140, 140}, 3)
	testResultLength(t, tree, [2]float32{15, 15}, [2]float32{135, 135}, 3)
	testResultLength(t, tree, [2]float32{0, 0}, [2]float32{120, 120}, 2)
	testResultLength(t, tree, [2]float32{60, 60}, [2]float32{120, 130}, 2)
	testResultLength(t, tree, [2]float32{150, 150}, [2]float32{160, 160}, 0)
}

func testResultLength(t *testing.T, tree RTree, min, max [2]float32, length int) {
	var results []interface{}
	tree.Search(min, max, func(min, max [2]float32, value interface{}) bool {
		s, _ := value.(string)
		results = append(results, s)
		return true
	})
	if len(results) != length {
		t.Error(fmt.Sprintf("results for %v %v, expected %d, result %d", min, max, length, len(results)))
	}
}
