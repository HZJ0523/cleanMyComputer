package ui

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/platform/windows"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func (a *App) newDashboard() fyne.CanvasObject {
	title := widget.NewLabelWithStyle(i18n.T("app.title"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	var diskInfo string
	homedir, _ := os.UserHomeDir()
	platform := windows.NewPlatform()
	if usage, err := platform.GetDiskUsage(homedir); err == nil {
		diskInfo = fmt.Sprintf(i18n.T("label.disk_usage"), usage.UsedGB, usage.TotalGB, usage.FreeGB)
	}

	diskLabel := widget.NewLabel(diskInfo)
	statusLabel := widget.NewLabel(i18n.T("label.select_mode"))
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	quickScan := widget.NewButton(i18n.T("btn.quick_scan"), func() {
		statusLabel.SetText(i18n.T("label.scanning"))
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
				statusLabel.SetText(fmt.Sprintf("%s, %d", i18n.T("label.scan_complete"), count))
				progressBar.Hide()
				dialog.ShowInformation(i18n.T("label.scan_complete"),
					fmt.Sprintf("%d items", count), a.window)
				a.tabs.SelectIndex(1)
			})
		}()
	})

	fullScan := widget.NewButton(i18n.T("btn.full_scan"), func() {
		statusLabel.SetText(i18n.T("label.scanning"))
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
				statusLabel.SetText(fmt.Sprintf("%s, %d", i18n.T("label.scan_complete"), count))
				progressBar.Hide()
				dialog.ShowInformation(i18n.T("label.scan_complete"),
					fmt.Sprintf("%d items", count), a.window)
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
