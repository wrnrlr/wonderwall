package colorpicker

import "gioui.org/widget"

func newHexField(editor widget.Editor, value string) *hexField {
	newValue, ok := parseHex(value)
	he := &hexField{Editor: editor, old: newValue}
	if ok {
		he.SetText(value)
	}
	he.Changed()
	return he
}

type hexField struct {
	widget.Editor
	Invalid bool
	old     []byte
}

func (ed *hexField) Changed() bool {
	s := ed.Editor.Text()
	newValue, ok := parseHex(s)
	if !ok {
		return false
	}
	changed := bytesEqual(newValue, ed.old)
	ed.old = newValue
	return changed
}

// SetText sets editor content without marking the editor changed.
func (ed *hexField) SetText(s string) {
	newValue, ok := parseHex(s)
	if !ok {
		return
	}
	ed.old = newValue
	ed.Editor.SetText(s)
}

type byteField struct {
	widget.Editor
	Invalid bool
	old     byte
}

func (ed *byteField) Changed() bool {
	newText := ed.Editor.Text()
	newByte, ok := parseByte(newText)
	if !ok {
		return false
	}
	changed := newByte != ed.old
	ed.old = newByte
	return changed
}

func (ed *byteField) SetText(s string) {
	newByte, ok := parseByte(s)
	if !ok {
		return
	}
	ed.old = newByte
	ed.Editor.SetText(s)
}

type percentageField struct {
	widget.Editor
	Invalid bool
	old     int
}

func (ed *percentageField) Changed() bool {
	newText := ed.Editor.Text()
	newInt, ok := parsePercentage(newText)
	if !ok {
		return false
	}
	changed := newInt != ed.old
	ed.old = newInt
	return changed
}

// SetText sets editor content without marking the editor changed.
func (ed *percentageField) SetText(s string) {
	newInt, ok := parsePercentage(s)
	if !ok {
		return
	}
	ed.old = newInt
	ed.Editor.SetText(s)
}

type degreeField struct {
	widget.Editor
	Invalid bool

	old int
}

func (ed *degreeField) Changed() bool {
	newText := ed.Editor.Text()
	i, ok := parseDegree(newText)
	if !ok {
		return false
	}
	degree := int(i)
	changed := degree != ed.old
	ed.old = degree
	return changed
}

// SetText sets editor content without marking the editor changed.
func (ed *degreeField) SetText(s string) {
	degree, ok := parseDegree(s)
	if !ok {
		return
	}
	ed.old = degree
	ed.Editor.SetText(s)
}
