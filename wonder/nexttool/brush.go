package nexttool

import "image/color"

type BrushOptions struct {
	StrokeColor color.NRGBA
}

type StrokeOptions struct {
	StrokeColor color.NRGBA
}

type FillOptions struct {
	Color         color.NRGBA
	GradietColor1 color.NRGBA
	GradietColor2 color.NRGBA
}
