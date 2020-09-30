package main

type Tool int

const (
	NoTool Tool = iota
	SelectionTool
	PenTool
	TextTool
	ImageTool
)

func (t Tool) String() string {
	switch t {
	case NoTool:
		return "NoTool"
	case SelectionTool:
		return "SelectionTool"
	case PenTool:
		return "PenTool"
	case TextTool:
		return "TextTool"
	default:
		return "UnknownTool"
	}
}
