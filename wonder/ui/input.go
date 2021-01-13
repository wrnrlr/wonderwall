package ui

import (
	"fmt"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image/color"
	"strconv"
)

func InputNumber(th *Theme, val int) *InputNumberStyle {
	editor := &widget.Editor{SingleLine: true}
	editor.SetText(fmt.Sprintf("%d", val))
	return &InputNumberStyle{
		Value:    val,
		Editor:   editor,
		TextSize: th.TextSize,
		Color:    th.Color.Text,
		shaper:   th.Shaper,
	}
}

type InputNumberStyle struct {
	Value    int
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.§
	Color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	Editor *widget.Editor

	shaper text.Shaper
}

func (in *InputNumberStyle) Layout(gtx C) D {
	defer op.Push(gtx.Ops).Pop()
	in.Event(gtx)
	dims := in.Editor.Layout(gtx, in.shaper, in.Font, in.TextSize)
	disabled := gtx.Queue == nil
	if in.Editor.Len() > 0 {
		textColor := in.Color
		if disabled {
			textColor = mulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(gtx.Ops)
		in.Editor.PaintText(gtx)
	}
	if !disabled {
		paint.ColorOp{Color: in.Color}.Add(gtx.Ops)
		in.Editor.PaintCaret(gtx)
	}
	return dims
}

func (in *InputNumberStyle) Event(gtx C) {
	for _, e := range in.Editor.Events() {
		_, ok := e.(widget.ChangeEvent)
		if !ok {
			return
		}
		txt := in.Editor.Text()
		if txt == "" {
			return
		}
		i, err := strconv.ParseInt(txt, 10, 64)
		if err != nil {
			in.Editor.SetText(fmt.Sprintf("%d", in.Value))
			return
		}
		num := int(i)
		in.Value = num
	}
}

func InputText(th *Theme, editor *widget.Editor, hint string) *Input {
	return &Input{
		Editor:    editor,
		TextSize:  th.TextSize,
		Color:     th.Color.Text,
		Hint:      hint,
		shaper:    th.Shaper,
		HintColor: th.Color.Hint,
	}
}

type Input struct {
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.
	Color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	Hint      string
	HintColor color.NRGBA
	Editor    *widget.Editor

	shaper text.Shaper
}

func (e Input) Layout(gtx C) D {
	defer op.Push(gtx.Ops).Pop()
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment}
	dims := tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = e.Editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	disabled := gtx.Queue == nil
	if e.Editor.Len() > 0 {
		textColor := e.Color
		if disabled {
			textColor = mulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(gtx.Ops)
		e.Editor.PaintText(gtx)
	} else {
		call.Add(gtx.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		e.Editor.PaintCaret(gtx)
	}
	return dims
}

func (in *Input) Event(gtx C) {

}
