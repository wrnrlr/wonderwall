package main

import (
	"gioui.org/io/pointer"
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

func (s *ImageService) Event(e pointer.Event, gtx C) interface{} {
	var result interface{}
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pos := e.Position.Mul(1 / gtx.Metric.PxPerDp)
	switch e.Type {
	case pointer.Press:
		src, err := s.Get("")
		if err != nil {
			break
		}
		result = AddImageEvent{Position: pos, Image: src}
	}
	return result
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
