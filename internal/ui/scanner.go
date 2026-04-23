package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewScannerView() fyne.CanvasObject {
	resultList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {},
	)
	return container.NewBorder(
		widget.NewLabel("扫描结果"),
		widget.NewButton("开始清理", func() {}),
		nil, nil,
		resultList,
	)
}
