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
		children = append(children, t.Icon(gtx, unit.Dp(36)))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func (m *ToolMenu) layoutOptions(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

type Arrange struct {
	btn widget.Clickable
}

func (t *Arrange) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, mouseIcon, sz)
}

func (t *Arrange) Event() {}

type Selection struct {
	btn widget.Clickable
}

func (t *Selection) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, selectionIcon, sz)
}

func (t *Selection) Event() {}

type Brush struct {
	btn widget.Clickable
}

func (t *Brush) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, brushIcon, sz)
}

func (t *Brush) Event() {}

type Pen struct {
	btn widget.Clickable
}

func (t *Pen) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, penIcon, sz)
}

func (t *Pen) Event() {}

type Text struct {
	btn widget.Clickable
}

func (t *Text) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, textIcon, sz)
}

func (t *Text) Event() {}

type Image struct {
	btn widget.Clickable
}

func (t *Image) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, imageIcon, sz)
}

func (t *Image) Event() {}

type Shape struct {
	btn widget.Clickable
}

func (t *Shape) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, shapeIcon, sz)
}

func (t *Shape) Event() {}

type Zoom struct {
	btn widget.Clickable
}

func (t *Zoom) Icon(gtx layout.Context, sz unit.Value) layout.FlexChild {
	return menuItem(&t.btn, zoomIcon, sz)
}

func (t *Zoom) Event() {}

func menuItem(btn *widget.Clickable, ic *widget.Icon, sz unit.Value) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
			return ic.Layout(gtx, sz)
		})
	})
}
