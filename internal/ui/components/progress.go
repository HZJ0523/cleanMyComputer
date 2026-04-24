package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewLabeledProgress(labelText string) fyne.CanvasObject {
	label := widget.NewLabel(labelText)
	progress := widget.NewProgressBar()
	return container.NewVBox(label, progress)
}
