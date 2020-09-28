package shape

import (
	"gioui.org/f32"
	"math"
)

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
