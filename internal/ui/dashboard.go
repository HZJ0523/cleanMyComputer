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

	diskLabel := widget.NewLabel("")
	go func() {
		homedir, _ := os.UserHomeDir()
		platform := windows.NewPlatform()
		if usage, err := platform.GetDiskUsage(homedir); err == nil {
			text := fmt.Sprintf(i18n.T("label.disk_usage"), usage.UsedGB, usage.TotalGB, usage.FreeGB)
			fyne.Do(func() { diskLabel.SetText(text) })
		} else {
			fyne.Do(func() { diskLabel.SetText("") })
		}
	}()

	statusLabel := widget.NewLabel(i18n.T("label.select_mode"))
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	doScan := func(level int) {
		statusLabel.SetText(i18n.T("label.scanning"))
		progressBar.Show()
		go func() {
			err := a.state.RunScan(level)
			fyne.Do(func() {
				progressBar.Hide()
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}
				count := a.state.GetScanItemCount()
				statusLabel.SetText(fmt.Sprintf("%s, %d", i18n.T("label.scan_complete"), count))
				dialog.ShowInformation(i18n.T("label.scan_complete"),
					fmt.Sprintf("%d items", count), a.window)
				a.scannerPage.refreshData()
				a.selectTab(1)
			})
		}()
	}

	quickScan := widget.NewButton(i18n.T("btn.quick_scan"), func() { doScan(1) })
	deepScan := widget.NewButton(i18n.T("btn.deep_scan"), func() { doScan(2) })
	fullScan := widget.NewButton(i18n.T("btn.full_scan"), func() { doScan(3) })

	return container.NewVBox(
		title,
		diskLabel,
		statusLabel,
		progressBar,
		container.NewHBox(quickScan, deepScan, fullScan),
	)
}
