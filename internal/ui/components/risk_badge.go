package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func NewRiskBadge(score int) fyne.CanvasObject {
	var text string
	var bg color.Color
	switch {
	case score > 60:
		text = i18n.T("risk.high")
		bg = color.RGBA{R: 230, G: 76, B: 76, A: 255}
	case score > 30:
		text = i18n.T("risk.moderate")
		bg = color.RGBA{R: 230, G: 204, B: 51, A: 255}
	default:
		text = i18n.T("risk.safe")
		bg = color.RGBA{R: 76, G: 204, B: 102, A: 255}
	}

	label := widget.NewLabelWithStyle(fmt.Sprintf("%s (%d)", text, score), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	rect := canvas.NewRectangle(bg)
	rect.SetMinSize(fyne.NewSize(100, 28))
	return container.NewStack(rect, label)
}
