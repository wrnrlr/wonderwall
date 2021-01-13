package ui

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
)

type Theme struct {
	Shaper text.Shaper
	Color  struct {
		Primary color.NRGBA
		Light   color.NRGBA
		Dark    color.NRGBA
		Text    color.NRGBA
		Hint    color.NRGBA
		InvText color.NRGBA
	}
	TextSize unit.Value
}

func MenuTheme(fontCollection []text.FontFace) *Theme {
	t := &Theme{Shaper: text.NewCache(fontCollection)}
	t.Color.Primary = Rgb(0xeeeeee)
	t.Color.Light = Rgb(0xcccccc)
	t.Color.Text = Rgb(0x000000)
	t.Color.Hint = Rgb(0x9b9b9b)
	t.Color.InvText = Rgb(0xffffff)
	t.TextSize = unit.Sp(16)

	return t
}

func loadIcon(b []byte) *widget.Icon {
	icon, err := widget.NewIcon(b)
	if err != nil {
		panic(err)
	}
	return icon
}
