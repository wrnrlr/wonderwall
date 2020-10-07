package main

import (
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/paint"
	"github.com/Almanax/wonderwall/wonder/shape"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type ImageService struct {
	cache map[string]image.Image
}

func (s *ImageService) Get(location string) (image.Image, error) {
	img, nil := getimage("yoda.jpg")
	return img, nil
}

func (s *ImageService) Event(gtx C) shape.Shape {
	var sh shape.Shape
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	for _, e := range gtx.Events(s) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Type {
		case pointer.Press:
			src, err := s.Get("")
			if err != nil {
				break
			}
			scale := 1 / gtx.Metric.PxPerDp
			img := paint.NewImageOp(src)
			pos := e.Position.Mul(scale)
			sh = &shape.Image{X: pos.X, Y: pos.Y, Image: img}
		case pointer.Release, pointer.Cancel:
		case pointer.Drag:
		}
	}
	pointer.InputOp{Tag: s, Grab: false, Types: pointer.Press | pointer.Drag | pointer.Release}.Add(gtx.Ops)
	return sh
}

func getimage(s string) (image.Image, error) {
	i, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(i)
	if err != nil {
		return nil, err
	}
	return im, nil
}
