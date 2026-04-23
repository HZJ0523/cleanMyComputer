package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewDashboard() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("电脑垃圾清理工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("选择扫描模式开始清理")
	quickScan := widget.NewButton("快速扫描", func() {})
	fullScan := widget.NewButton("完整扫描", func() {})
	return container.NewVBox(
		title,
		subtitle,
		container.NewHBox(quickScan, fullScan),
	)
}
