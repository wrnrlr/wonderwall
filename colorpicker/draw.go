package colorpicker

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
	"math"
)

const c = 0.55228475 // 4*(sqrt(2)-1)/3

func drawControl(p f32.Point, radius, width float32, gtx layout.Context) {
	width = float32(gtx.Px(unit.Dp(width))) / 2
	radius = float32(gtx.Px(unit.Dp(radius))) - width
	p.X -= radius - width*4
	p.Y -= radius - width*9
	drawCircle(p, radius, width, color.NRGBA{A: 0xff}, gtx)
	p.X += width * 2
	p.Y += width * 2
	drawCircle(p, radius, width, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, gtx)
}

func drawCircle(p f32.Point, radius, width float32, col color.NRGBA, gtx layout.Context) {
	w, h := radius*2, radius*2
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Point{X: p.X, Y: p.Y})
	path.Move(f32.Point{X: w / 4 * 3, Y: radius / 2})
	path.Cube(f32.Point{X: 0, Y: radius * c}, f32.Point{X: -radius + radius*c, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * c, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*c}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * c}, f32.Point{X: radius - radius*c, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * c, Y: 0}, f32.Point{X: radius, Y: radius - radius*c}, f32.Point{X: radius, Y: radius})      // NE
	path.Move(f32.Point{X: -w, Y: -radius})                                                                                     // Return to origin
	scale := (radius - width*2) / radius
	path.Move(f32.Point{X: w * (1 - scale) * .5, Y: h * (1 - scale) * .5})
	w *= scale
	h *= scale
	radius *= scale
	path.Move(f32.Point{X: 0, Y: h - radius})
	path.Cube(f32.Point{X: 0, Y: radius * c}, f32.Point{X: +radius - radius*c, Y: radius}, f32.Point{X: +radius, Y: radius})      // SW
	path.Cube(f32.Point{X: +radius * c, Y: 0}, f32.Point{X: +radius, Y: -radius + radius*c}, f32.Point{X: +radius, Y: -radius})   // SE
	path.Cube(f32.Point{X: 0, Y: -radius * c}, f32.Point{X: -(radius - radius*c), Y: -radius}, f32.Point{X: -radius, Y: -radius}) // NE
	path.Cube(f32.Point{X: -radius * c, Y: 0}, f32.Point{X: -radius, Y: radius - radius*c}, f32.Point{X: -radius, Y: radius})     // NW
	clip.Outline{Path: path.End()}.Op().Add(gtx.Ops)
	cons := gtx.Constraints
	dr := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: cons.Max.X, Y: cons.Max.Y}}
	clip.Rect(dr).Add(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawCheckerboard(gtx layout.Context) {
	w := gtx.Constraints.Max.X
	h := gtx.Constraints.Max.Y
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, clip.Rect{Max: gtx.Constraints.Max}.Op())
	size := h / 2
	defer op.Save(gtx.Ops).Load()
	var path clip.Path
	path.Begin(gtx.Ops)
	count := int(math.Ceil(float64(w / size)))
	for i := 0; i < count; i++ {
		offset := 0
		if math.Mod(float64(i), 2) == 0 {
			offset += size
		}
		path.MoveTo(f32.Point{X: float32(i * size), Y: float32(offset)})
		path.Line(f32.Point{X: float32(size)})
		path.Line(f32.Point{Y: float32(size)})
		path.Line(f32.Point{X: float32(-size)})
		path.Line(f32.Point{Y: float32(-size)})
	}
	clip.Outline{Path: path.End()}.Op().Add(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0xef, G: 0xef, B: 0xef, A: 0xff}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
