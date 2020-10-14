package main

import (
	"gioui.org/f32"
	"github.com/Almanax/wonderwall/wonder/shape"
	"github.com/rs/xid"
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

type InsertShapeEvent struct{}

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
