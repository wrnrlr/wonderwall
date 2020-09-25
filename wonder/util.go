package main

import (
	"gioui.org/f32"
	"gioui.org/widget"
	"image"
	"math"
)

const (
	rad90 = float32(90 * math.Pi / 180)
)

func loadIcon(b []byte) *widget.Icon {
	icon, err := widget.NewIcon(b)
	if err != nil {
		panic(err)
	}
	return icon
}

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

func offsetPoint(point f32.Point, distance, angle float32) f32.Point {
	x := point.X + distance*cos(angle)
	y := point.Y + distance*sin(angle)
	return f32.Point{X: x, Y: y}
}

func angle(p1, p2 f32.Point) float32 {
	return atan2(p2.Y-p1.Y, p2.X-p1.X)
}

func cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func imageRect(r f32.Rectangle) image.Rectangle {
	x1, y1 := int(r.Min.X), int(r.Min.Y)
	x2, y2 := int(r.Max.X), int(r.Max.Y)
	return image.Rect(x1, y1, x2, y2)
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}
