package nexttool

import (
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func loadIcon(b []byte) *widget.Icon {
	icon, err := widget.NewIcon(b)
	if err != nil {
		panic(err)
	}
	return icon
}

var (
	mouseIcon     = loadIcon(icons.ActionPanTool)
	brushIcon     = loadIcon(icons.ImageBrush)
	textIcon      = loadIcon(icons.EditorTitle)
	deleteIcon    = loadIcon(icons.ActionDelete)
	undoIcon      = loadIcon(icons.ContentUndo)
	redoIcon      = loadIcon(icons.ContentRedo)
	backIcon      = loadIcon(icons.NavigationArrowBack)
	filterIcon    = loadIcon(icons.HardwareKeyboardArrowDown)
	sortIcon      = loadIcon(icons.ContentSort)
	addBoxIcon    = loadIcon(icons.ContentAddBox)
	userIcon      = loadIcon(icons.SocialPerson)
	imageIcon     = loadIcon(icons.ImageImage)
	selectionIcon = loadIcon(icons.ImageCropFree)
	penIcon       = loadIcon(icons.ImageEdit)
	shapeIcon     = loadIcon(icons.ImageVignette)
	zoomIcon      = loadIcon(icons.ActionSearch)
	shape2Icon    = loadIcon(icons.ImageLandscape)
	shape3Icon    = loadIcon(icons.MapsStreetView)
)
