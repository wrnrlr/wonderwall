package colorpicker

import "image/color"

type ColorInput interface {
	Changed() bool
	SetColor(col color.NRGBA)
	Color() color.NRGBA
}
