package ui

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
	"image/color"
)

type ColorPicker struct {
	theme      *Theme
	Button     *widget.Clickable
	Background color.RGBA
	// Color is the icon color.
	Size  unit.Value
	Inset layout.Inset

	Color color.RGBA

	hex *widget.Editor

	active bool
}

func Color(th *Theme) *ColorPicker {
	return &ColorPicker{
		theme:  th,
		Color:  th.Color.InvText,
		Size:   unit.Dp(24),
		Inset:  layout.UniformInset(unit.Dp(12)),
		Button: &widget.Clickable{},
		hex:    &widget.Editor{SingleLine: true},
		active: false,
	}
}

func (cp *ColorPicker) Layout(gtx C) D {
	dims := cp.layoutButton(gtx)
	if cp.active {
		stack := op.Push(gtx.Ops)
		c := gtx.Constraints
		op.Offset(f32.Pt(0, float32(dims.Size.Y))).Add(gtx.Ops)
		cp.layoutPanel(gtx)
		gtx.Constraints = c
		stack.Pop()
	}
	return dims
}

func (cp *ColorPicker) layoutButton(gtx C) D {
	width := int(unit.Dp(40).V)
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			clip.Rect{Max: gtx.Constraints.Min}.Add(gtx.Ops)
			background := cp.Background
			if gtx.Queue == nil {
				background = mulAlpha(cp.Background, 150)
			}
			dims := Fill(gtx, background)
			for _, c := range cp.Button.History() {
				drawInk(gtx, c)
			}
			return dims
		}),
		layout.Stacked(func(gtx C) D {
			return cp.Inset.Layout(gtx, func(gtx C) D {
				height := int(cp.Size.V)
				gtx.Constraints.Min.X = width
				gtx.Constraints.Min.Y = height
				Fill(gtx, Rgb(0x0000ff))
				return D{Size: image.Point{X: width, Y: height}}
			})
		}),
		layout.Expanded(func(gtx C) D {
			pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return cp.Button.Layout(gtx)
		}),
	)
}

func (cp *ColorPicker) layoutPanel(gtx C) D {
	gtx.Constraints.Max.X = 400
	colors := layout.Rigid(func(gtx C) D {
		grid := Grid{Columns: 6, Rows: 5, Width: 600, Height: 400}
		return grid.Layout(gtx, func(i, j int, gtx C) D {
			col := colorPalet[i*grid.Columns+j]
			return Fill(gtx, Rgb(col))
		})
	})
	hexinput := layout.Rigid(func(gtx C) D {
		return Editor(cp.theme, cp.hex, "#hexval").Layout(gtx)
	})
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, colors, hexinput)
}

func (cp *ColorPicker) Event(gtx C) {
	if cp.Button.Clicked() {
		cp.active = !cp.active
	}
}

var colorPalet []uint32 = []uint32{
	0x41B0F6,
	0x74FBEA,
	0x89F94F,
	0xFFFC67,
	0xFE968D,
	0xFF8EC6,

	0x00A1FE,
	0x1EE5CE,
	0x60D838,
	0xF9E231,
	0xFF634D,
	0xEE5FA7,

	0x0376BB,
	0x05A89D,
	0x1EB100,
	0xF7BA00,
	0xED230D,
	0xCA2A7A,

	0x004D81,
	0x007C77,
	0x047101,
	0xFF9400,
	0xB51700,
	0x9A1860,

	0xFFFFFF,
	0xD5D5D5,
	0x929292,
	0x646464,
	0x000000,
	0xFFC943,
}
