package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewHistoryView() fyne.CanvasObject {
	historyList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)
	return container.NewBorder(
		widget.NewLabel("清理历史"),
		nil, nil, nil,
		historyList,
	)
}
