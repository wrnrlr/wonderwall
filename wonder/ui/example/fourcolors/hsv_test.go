package main

import (
	"fmt"
	"image/color"
	"testing"
)

func TestRgbHsvConversion(t *testing.T) {
	rgbHsvs := []struct {
		rgb color.RGBA
		hsv HSVColor
	}{
		{color.RGBA{A: 0xff}, HSVColor{H: 0, S: 0, V: 0}},                              // black
		{color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 0, S: 0, V: 1}},   // white
		{color.RGBA{R: 0xff, A: 0xff}, HSVColor{H: 0, S: 1, V: 1}},                     // red
		{color.RGBA{G: 0xff, A: 0xff}, HSVColor{H: 120, S: 1, V: 1}},                   // lime
		{color.RGBA{B: 0xff, A: 0xff}, HSVColor{H: 240, S: 1, V: 1}},                   // blue
		{color.RGBA{R: 0xff, G: 0xff, A: 0xff}, HSVColor{H: 60, S: 1, V: 1}},           // yellow
		{color.RGBA{G: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 180, S: 1, V: 1}},          // cyan
		{color.RGBA{R: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 300, S: 1, V: 1}},          // magenta
		{color.RGBA{R: 0xbf, G: 0xbf, B: 0xbf, A: 0xff}, HSVColor{H: 0, S: 1, V: .75}}, // silver
		{color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 0, S: 1, V: .5}},  // gray
		{color.RGBA{R: 0x80, A: 0xff}, HSVColor{H: 0, S: 1, V: .5}},                    // maroon
		{color.RGBA{R: 0xff, G: 0x80, A: 0xff}, HSVColor{H: 60, S: 1, V: .5}},          // olive
		{color.RGBA{G: 0x80, A: 0xff}, HSVColor{H: 120, S: 1, V: .5}},                  // green
		{color.RGBA{R: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 300, S: 1, V: .5}},         // purple
		{color.RGBA{G: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 180, S: 1, V: .5}},         // teal
		{color.RGBA{B: 0x80, A: 0xff}, HSVColor{H: 240, S: 1, V: .5}},                  // navy
	}
	for _, v := range rgbHsvs {
		testHsv(v.hsv, RgbToHsv(v.rgb))
		testRgb(v.rgb, HsvToRgb(v.hsv))
	}
}

func testRgb(a color.RGBA, b color.RGBA) {
	if a.R != b.R && a.G != b.G && a.B != b.B {
		panic(fmt.Sprintf("not equal: %v, %v", a, b))
	}
}

func testHsv(a HSVColor, b HSVColor) {
	if a.H != b.H && a.S != b.S && a.V != b.V {
		panic(fmt.Sprintf("not equal: %v, %v", a, b))
	}
}
