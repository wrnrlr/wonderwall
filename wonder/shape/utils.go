package shape

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"image"
	"math"
)

type (
	D = layout.Dimensions
	C = layout.Context
)

const (
	rad90 = float32(90 * math.Pi / 180)
)

func boundingBox(points []f32.Point) (box f32.Rectangle) {
	box.Min, box.Max = points[0], points[0]
	for _, p := range points {
		box.Min.X = min(box.Min.X, p.X)
		box.Min.Y = min(box.Min.Y, p.Y)
		box.Max.X = max(box.Max.X, p.X)
		box.Max.Y = max(box.Max.Y, p.Y)
	}
	return box
}

func circleBox(p f32.Point, r float32) f32.Rectangle {
	return f32.Rect(p.X-r, p.Y-r, p.X+r, p.Y+r)
}

func offsetPoint(point f32.Point, distance, angle float32) f32.Point {
	x := point.X + distance*cos(angle)
	y := point.Y + distance*sin(angle)
	return f32.Point{X: x, Y: y}
}

func imageRect(r f32.Rectangle) image.Rectangle {
	x1, y1 := int(r.Min.X), int(r.Min.Y)
	x2, y2 := int(r.Max.X), int(r.Max.Y)
	return image.Rect(x1, y1, x2, y2)
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rect(float32(r.Min.X), float32(r.Min.Y), float32(r.Max.X), float32(r.Max.Y))
}
