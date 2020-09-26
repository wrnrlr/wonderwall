package main

import "github.com/rs/xid"

type BackEvent struct{}

type LoginEvent struct{}

type ShowWallListEvent struct{}

type ShowWallEvent struct {
	WallID xid.ID
}

type ShowAddWallEvent struct{}

type ShowUserEvent struct{}
