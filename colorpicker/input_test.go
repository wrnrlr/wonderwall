package colorpicker

import (
	"image/color"
	"testing"
)

func TestInputSetColor(t *testing.T) {
	red := color.NRGBA{R: 0xff, A: 0xff}
	inputs := []ColorInput{
		NewPicker(nil),
		NewAlphaSlider(),
		NewHexEditor(nil),
		NewRgbEditor(nil),
		NewHsvEditor(nil),
		NewToggle(nil, NewHexEditor(nil)),
		NewColorSelection(nil),
		NewMux(),
	}
	for _, in := range inputs {
		if in.Changed() {
			t.Fail()
		}
		in.SetColor(red)
		if in.Changed() {
			t.Fail()
		}
		if in.Color() != red {
			t.Fail()
		}
	}
}
