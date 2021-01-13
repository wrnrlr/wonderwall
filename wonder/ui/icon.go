package ui

import (
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

type Icon struct {
	Src  []byte
	Size unit.Value

	// Cached values.
	op      paint.ImageOp
	imgSize int
}

func (ic *Icon) Image(c unit.Metric, col color.NRGBA) paint.ImageOp {
	sz := c.Px(ic.Size)
	if sz == ic.imgSize {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.Src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	r, g, b, a := col.RGBA()
	m.Palette[0] = color.RGBA{r, rgba.G, rgba.B, a}
	iconvg.Decode(&ico, ic.Src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	return ic.op
}
