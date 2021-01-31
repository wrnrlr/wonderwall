package main

import (
	"fmt"
	"image/color"
	"testing"
)

func TestRgbToHsv(t *testing.T) {
	equal(color.NRGBA{A: 0xff}, HSVColor{H: 0, S: 0, V: 0})                              // black
	equal(color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 0, S: 0, V: 1})   // white
	equal(color.NRGBA{R: 0xff, A: 0xff}, HSVColor{H: 0, S: 1, V: 1})                     // red
	equal(color.NRGBA{G: 0xff, A: 0xff}, HSVColor{H: 120, S: 1, V: 1})                   // lime
	equal(color.NRGBA{B: 0xff, A: 0xff}, HSVColor{H: 240, S: 1, V: 1})                   // blue
	equal(color.NRGBA{R: 0xff, G: 0xff, A: 0xff}, HSVColor{H: 60, S: 1, V: 1})           // yellow
	equal(color.NRGBA{G: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 180, S: 1, V: 1})          // cyan
	equal(color.NRGBA{R: 0xff, B: 0xff, A: 0xff}, HSVColor{H: 300, S: 1, V: 1})          // magenta
	equal(color.NRGBA{R: 0xbf, G: 0xbf, B: 0xbf, A: 0xff}, HSVColor{H: 0, S: 1, V: .75}) // silver
	equal(color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 0, S: 1, V: .5})  // gray
	equal(color.NRGBA{R: 0x80, A: 0xff}, HSVColor{H: 0, S: 1, V: .5})                    // maroon
	equal(color.NRGBA{R: 0xff, G: 0x80, A: 0xff}, HSVColor{H: 60, S: 1, V: .5})          // olive
	equal(color.NRGBA{G: 0x80, A: 0xff}, HSVColor{H: 120, S: 1, V: .5})                  // green
	equal(color.NRGBA{R: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 300, S: 1, V: .5})         // purple
	equal(color.NRGBA{G: 0x80, B: 0x80, A: 0xff}, HSVColor{H: 180, S: 1, V: .5})         // teal
	equal(color.NRGBA{B: 0x80, A: 0xff}, HSVColor{H: 240, S: 1, V: .5})                  // navy
}

func equal(rgb color.NRGBA, b HSVColor) {
	a := RgbToHsv(rgb)
	if a.H != b.H && a.S != b.S && a.V != b.V {
		panic(fmt.Sprintf("not equal: %v, %v", a, b))
	}
}
