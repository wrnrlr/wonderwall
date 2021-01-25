package ui

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/wrnrlr/wonderwall/wonder/f32color"
	"image"
	"image/color"
)

type ColorPicker struct {
	theme      *Theme
	Button     *widget.Clickable
	Background color.NRGBA
	// Color is the icon color.
	Size  unit.Value
	Inset layout.Inset

	Color color.NRGBA

	grid    *Grid
	buttons [30]widget.Clickable

	hex *widget.Editor

	active bool
}

func Color(th *Theme, col color.NRGBA) *ColorPicker {
	return &ColorPicker{
		theme:   th,
		Color:   col,
		Size:    unit.Dp(24),
		Inset:   layout.UniformInset(unit.Dp(12)),
		Button:  &widget.Clickable{},
		buttons: [30]widget.Clickable{},
		hex:     &widget.Editor{SingleLine: true},
		active:  false,
	}
}

func (cp *ColorPicker) Layout(gtx C) D {
	stack := op.Save(gtx.Ops)
	dims := cp.layoutButton(gtx)
	stack.Load()
	if cp.active {
		stack := op.Save(gtx.Ops)
		cons := gtx.Constraints
		op.Offset(f32.Pt(0, float32(dims.Size.Y))).Add(gtx.Ops)
		cp.layoutPanel(gtx)
		gtx.Constraints = cons
		stack.Load()
	}
	return dims
}

func (cp *ColorPicker) layoutButton(gtx C) D {
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
				height := gtx.Px(cp.Size)
				gtx.Constraints = layout.Exact(image.Point{X: height, Y: height})
				Fill(gtx, cp.Color)
				return D{Size: image.Point{X: height, Y: height}}
			})
		}),
		layout.Expanded(func(gtx C) D {
			pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return cp.Button.Layout(gtx)
		}),
	)
}

func (cp *ColorPicker) layoutPanel(gtx C) D {
	colors := layout.Rigid(func(gtx C) D {
		size := image.Point{X: 300, Y: 200}
		gtx.Constraints = layout.Constraints{Min: size, Max: size}
		return Grid{Columns: 6, Rows: 5}.Layout(gtx, func(i, j int, gtx C) {
			index := i*6 + j
			cp.buttons[index].Layout(gtx)
			Fill(gtx, f32color.RGBAToNRGBA(Rgb(ColorPalet[index])))
		})
	})
	hexinput := layout.Rigid(func(gtx C) D {
		return Editor(cp.theme, cp.hex, "#hexval").Layout(gtx)
	})
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, colors, hexinput)
}

func (cp *ColorPicker) Event(gtx C) (col *color.RGBA) {
	if cp.Button.Clicked() {
		cp.active = !cp.active
	}
	if !cp.active {
		return nil
	}
	for i := range cp.buttons {
		if cp.buttons[i].Clicked() {
			fmt.Printf("Color grid clicked: \n")
			cp.Color = f32color.RGBAToNRGBA(Rgb(ColorPalet[i]))
		}
	}
	return col
}

var ColorPalet = []uint32{
	0xaf0000, // lightblue
	0x74FBEA, //
	0x89F94F,
	0xFFFC67,
	0xFE968D,
	0xFF8EC6, // lightpink

	0x00A1FE, // blue
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
