package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/platform/windows"
)

func (a *App) newDashboard() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("电脑垃圾清理工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	var diskInfo string
	homedir, _ := os.UserHomeDir()
	platform := windows.NewPlatform()
	if usage, err := platform.GetDiskUsage(homedir); err == nil {
		diskInfo = fmt.Sprintf("系统盘: 已用 %.1f GB / 总计 %.1f GB (可用 %.1f GB)",
			usage.UsedGB, usage.TotalGB, usage.FreeGB)
	}

	diskLabel := widget.NewLabel(diskInfo)
	statusLabel := widget.NewLabel("选择扫描模式开始清理")
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	quickScan := widget.NewButton("快速扫描", func() {
		statusLabel.SetText("正在快速扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(1)
			if err != nil {
				fyne.Do(func() {
					progressBar.Hide()
					dialog.ShowError(err, a.window)
				})
				return
			}
			count := len(a.state.ScanItems)
			fyne.Do(func() {
				statusLabel.SetText(fmt.Sprintf("扫描完成，发现 %d 个可清理项", count))
				progressBar.Hide()
				dialog.ShowInformation("扫描完成",
					fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", count), a.window)
				a.tabs.SelectIndex(1)
			})
		}()
	})

	fullScan := widget.NewButton("完整扫描", func() {
		statusLabel.SetText("正在完整扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(3)
			if err != nil {
				fyne.Do(func() {
					progressBar.Hide()
					dialog.ShowError(err, a.window)
				})
				return
			}
			count := len(a.state.ScanItems)
			fyne.Do(func() {
				statusLabel.SetText(fmt.Sprintf("扫描完成，发现 %d 个可清理项", count))
				progressBar.Hide()
				dialog.ShowInformation("扫描完成",
					fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", count), a.window)
				a.tabs.SelectIndex(1)
			})
		}()
	})

	return container.NewVBox(
		title,
		diskLabel,
		statusLabel,
		progressBar,
		container.NewHBox(quickScan, fullScan),
	)
}
