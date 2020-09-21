package ui

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

type Theme struct {
	Shaper text.Shaper
	Color  struct {
		Primary color.RGBA
		Light   color.RGBA
		Dark    color.RGBA
		Text    color.RGBA
		Hint    color.RGBA
		InvText color.RGBA
	}
	TextSize unit.Value
}

func CustomTheme(fontCollection []text.FontFace) *material.Theme {
	t := &material.Theme{Shaper: text.NewCache(fontCollection)}
	t.Color.Primary = Rgb(0xcccccc)
	t.Color.Text = Rgb(0xffffff)
	t.Color.Hint = Rgb(0x9b9b9b)
	t.Color.InvText = Rgb(0x000000)
	t.TextSize = unit.Sp(16)
	return t
}

func MenuTheme(fontCollection []text.FontFace) *Theme {
	t := &Theme{Shaper: text.NewCache(fontCollection)}
	t.Color.Primary = Rgb(0xcccc)
	t.Color.Light = Rgb(0xcccc)
	t.Color.Primary = Rgb(0xcccc)
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
