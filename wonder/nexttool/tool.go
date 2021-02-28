package nexttool

import (
	"gioui.org/widget"
)

type Tool interface {
	Icon() *widget.Icon
	Event()
}
