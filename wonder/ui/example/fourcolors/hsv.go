// Implementation of HSV colorspace

package main

import "image/color"

type HSVColor struct {
	H uint8
	S uint8
	V uint8
}

func HsvToRgb(hsv HSVColor) color.RGBA {
	var rgb color.RGBA

	var region, p, q, t uint8
	var h, s, v, remainder uint16

	if hsv.S == 0 {
		rgb.R = hsv.V
		rgb.G = hsv.V
		rgb.B = hsv.V
		return rgb
	}

	// convert to 16bit to prevent overflow
	h = uint16(hsv.H)
	s = uint16(hsv.S)
	v = uint16(hsv.V)

	region = uint8(h / 43)
	remainder = (h - uint16(region*43)) * 6

	p = uint8((v * (255 - s)) >> 8)
	q = uint8((v * (255 - ((s * remainder) >> 8))) >> 8)
	t = uint8((v * (255 - ((255 - (255 - remainder)) >> 8))) >> 8)

	switch region {
	case 0:
		rgb.R = uint8(v)
		rgb.G = t
		rgb.B = p
	case 1:
		rgb.R = q
		rgb.G = uint8(v)
		rgb.B = p
	case 2:
		rgb.R = p
		rgb.G = uint8(v)
		rgb.B = t
	case 3:
		rgb.R = p
		rgb.G = q
		rgb.B = uint8(v)
	case 4:
		rgb.R = t
		rgb.G = p
		rgb.B = uint8(v)
	default:
		rgb.R = uint8(v)
		rgb.G = p
		rgb.B = q
	}
	return rgb
}

func RgbToHsv(rgb color.NRGBA) HSVColor {

	var hsv HSVColor

	rgbMin := min(min(rgb.R, rgb.G), rgb.B)
	rgbMax := max(max(rgb.R, rgb.G), rgb.B)

	hsv.V = rgbMax
	if hsv.V == 0 {
		return hsv
	}

	hsv.S = 255 * (rgbMax - rgbMin) / hsv.V
	if hsv.S == 0 {
		return hsv
	}

	if rgbMax == rgb.R {
		hsv.H = 0 + 43*(rgb.G-rgb.B)/(rgbMax-rgbMin)
	} else if rgbMax == rgb.G {
		hsv.H = 85 + 43*(rgb.B-rgb.R)/(rgbMax-rgbMin)
	} else {
		hsv.H = 171 + 43*(rgb.R-rgb.G)/(rgbMax-rgbMin)
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
