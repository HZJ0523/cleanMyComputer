package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newDashboard() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("电脑垃圾清理工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabel("选择扫描模式开始清理")

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	quickScan := widget.NewButton("快速扫描", func() {
		statusLabel.SetText("正在快速扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(1)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			statusLabel.SetText("扫描完成")
			progressBar.Hide()
			msg := fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", len(a.state.ScanItems))
			dialog.ShowInformation("扫描完成", msg, a.window)
			a.selectTab(1)
		}()
	})

	fullScan := widget.NewButton("完整扫描", func() {
		statusLabel.SetText("正在完整扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(3)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			statusLabel.SetText("扫描完成")
			progressBar.Hide()
			msg := fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", len(a.state.ScanItems))
			dialog.ShowInformation("扫描完成", msg, a.window)
			a.selectTab(1)
		}()
	})

	return container.NewVBox(
		title,
		statusLabel,
		progressBar,
		container.NewHBox(quickScan, fullScan),
	)
}
