package nexttool

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func NewToolMenu(theme *material.Theme, tools ...Tool) *ToolMenu {
	return &ToolMenu{theme, tools}
}

type ToolMenu struct {
	theme *material.Theme
	tools []Tool
}

func (m *ToolMenu) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(m.layoutOptions),
		layout.Rigid(m.layoutMenu))
}

func (m *ToolMenu) layoutMenu(gtx layout.Context) layout.Dimensions {
	var children []layout.FlexChild
	for _, t := range m.tools {
		ic := t.Icon()
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ic.Layout(gtx, unit.Dp(36))
			}))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (m *ToolMenu) layoutOptions(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

type Arrange struct {
	btn widget.Clickable
}

func (t *Arrange) Icon() *widget.Icon {
	return mouseIcon
}

func (t *Arrange) Event() {}

type Selection struct {
	btn widget.Clickable
}

func (t *Selection) Icon() *widget.Icon {
	return selectionIcon
}

func (t *Selection) Event() {}

type Brush struct {
	btn widget.Clickable
}

func (t *Brush) Icon() *widget.Icon {
	return brushIcon
}

func (t *Brush) Event() {}

type Pen struct {
	btn widget.Clickable
}

func (t *Pen) Icon() *widget.Icon {
	return penIcon
}

func (t *Pen) Event() {}

type Text struct {
	btn widget.Clickable
}

func (t *Text) Icon() *widget.Icon {
	return textIcon
}

func (t *Text) Event() {}

type Image struct {
	btn widget.Clickable
}

func (t *Image) Icon() *widget.Icon {
	return imageIcon
}

func (t *Image) Event() {}

type Shape struct {
	btn widget.Clickable
}

func (t *Shape) Icon() *widget.Icon {
	return shapeIcon
}

func (t *Shape) Event() {}

type Zoom struct {
	btn widget.Clickable
}

func (t *Zoom) Icon() *widget.Icon {
	return zoomIcon
}

func (t *Zoom) Event() {}
