package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewConfirmView() fyne.CanvasObject {
	summary := widget.NewLabel("准备清理 0 个文件，释放 0 B 空间")
	confirmBtn := widget.NewButton("确认清理", func() {})
	cancelBtn := widget.NewButton("取消", func() {})
	return container.NewVBox(summary, container.NewHBox(confirmBtn, cancelBtn))
}
