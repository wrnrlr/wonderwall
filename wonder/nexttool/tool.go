package nexttool

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type Tool interface {
	Icon(gtx layout.Context, sz unit.Value) layout.FlexChild
	Event()
}
