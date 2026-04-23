package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsView() fyne.CanvasObject {
	autoClean := widget.NewCheck("启用自动清理", func(bool) {})
	retentionEntry := widget.NewEntry()
	retentionEntry.SetPlaceHolder("24")
	return container.NewVBox(
		widget.NewLabel("设置"),
		autoClean,
		widget.NewLabel("隔离区保留时间（小时）"),
		retentionEntry,
	)
}
