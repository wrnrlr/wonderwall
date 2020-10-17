package main

import (
	"gioui.org/f32"
	"github.com/Almanax/wonderwall/wonder/shape"
	"github.com/rs/xid"
	"image"
)

type BackEvent struct{}

type LoginEvent struct{}

type ShowWallListEvent struct{}

type ShowWallEvent struct {
	WallID xid.ID
}

type ShowAddWallEvent struct{}

type ShowUserEvent struct{}

type DeleteEvent struct{}

type AddLineEvent struct {
	Points []f32.Point
}

type AddTextEvent struct {
	Position f32.Point
}

type AddImageEvent struct {
	Position f32.Point
	Image    image.Image
}

type InsertShapeEvent struct {
	Shape shape.Shape
}

type UpdateShapeEvent struct {
	Shape shape.Shape
}

type MoveShapeEvent struct {
	Offset f32.Point
}

type SelectionEvent struct {
	Shape shape.Shape
}

type ZoomEvent struct {
	Scroll float32
	Pos    f32.Point
}

type PanEvent struct {
	Offset f32.Point
	Pos    f32.Point
}
