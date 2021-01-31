// Implementation of HSV colorspace

package main

import (
	"image/color"
)

type HSVColor struct {
	H float32 // [0,1]
	S float32 // [0,1]
	V float32 // [0,360]
}

func HsvToRgb(hsv HSVColor) color.RGBA {
	v := uint8(hsv.V * 255)
	if hsv.S == 0.0 {
		return color.RGBA{v, v, v, 0xff}
	}

	hh := hsv.H
	if hh >= 360.0 {
		hh = 0.0
	}

	hh /= 60.0
	i := int(hh)
	ff := hh - float32(i)
	p := uint8(hsv.V * (1 - hsv.S) * 255)
	q := uint8(hsv.V * (1 - (hsv.S * ff)) * 255)
	t := uint8(hsv.V * (1 - (hsv.S * (1 - ff))) * 255)

	rgb := color.RGBA{A: 0xff}
	switch i {
	case 0:
		rgb.R = v
		rgb.G = t
		rgb.B = p
	case 1:
		rgb.R = q
		rgb.G = v
		rgb.B = p
	case 2:
		rgb.R = p
		rgb.G = v
		rgb.B = t
	case 3:
		rgb.R = p
		rgb.G = q
		rgb.B = v
	case 4:
		rgb.R = t
		rgb.G = p
		rgb.B = v
	case 5:
	default:
		rgb.R = v
		rgb.G = p
		rgb.B = q
	}
	return rgb
}

func RgbToHsv(rgb color.NRGBA) HSVColor {

	var hsv HSVColor

	rgbMin := min(min(rgb.R, rgb.G), rgb.B)
	rgbMax := max(max(rgb.R, rgb.G), rgb.B)

	hsv.V = float32(rgbMax / 255)
	delta := rgbMax - rgbMin
	if delta < 0 {
		return hsv
	}
	if rgbMax == 0 {
		//hsv.H = NaN
		return hsv
	} else {
		hsv.S = float32(delta / rgbMax) // s
	}
	if rgb.R >= rgbMax {
		hsv.H = float32((rgb.G - rgb.B) / delta)
	} else if rgb.R >= rgbMax {
		hsv.H = float32(2.0 + (rgb.B*255-rgb.R*255)/delta)
	} else {
		hsv.H = float32(4.0 + (rgb.R*255-rgb.G*255)/delta)
	}
	hsv.H *= 60
	if hsv.H < 0 {
		hsv.H += 360
	}
	return hsv
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	} else {
		return b
	}
}
