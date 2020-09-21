package ui

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type ColorPicker struct {
	Button     *widget.Clickable
	Background color.RGBA
	// Color is the icon color.
	Size  unit.Value
	Inset layout.Inset
	Color color.RGBA
}

func Color(th *material.Theme, button *widget.Clickable) ColorPicker {
	return ColorPicker{
		Background: th.Color.Primary,
		Color:      th.Color.InvText,
		Size:       unit.Dp(24),
		Inset:      layout.UniformInset(unit.Dp(12)),
		Button:     button,
	}
}

func (b ColorPicker) Layout(gtx C) D {
	width := int(unit.Dp(40).V)
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			clip.Rect{Max: gtx.Constraints.Min}.Add(gtx.Ops)

			background := b.Background
			if gtx.Queue == nil {
				background = mulAlpha(b.Background, 150)
			}
			dims := Fill(gtx, background)
			for _, c := range b.Button.History() {
				drawInk(gtx, c)
			}
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				height := int(b.Size.V)
				gtx.Constraints.Min.X = width
				gtx.Constraints.Min.Y = height
				Fill(gtx, Rgb(0x0000ff))
				return layout.Dimensions{Size: image.Point{X: width, Y: height}}
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return b.Button.Layout(gtx)
		}),
	)
}
