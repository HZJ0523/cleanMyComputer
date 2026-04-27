package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func (a *App) newConfirmView() fyne.CanvasObject {
	summaryLabel := widget.NewLabel(i18n.T("label.please_scan"))

	confirmBtn := widget.NewButton(i18n.T("btn.confirm_clean"), func() {
		items := a.state.GetScanItemsSafe()
		if len(items) == 0 {
			dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("label.no_items"), a.window)
			return
		}

		var highRiskCount int
		for _, item := range items {
			if item.RiskScore > 60 {
				highRiskCount++
			}
		}
		if highRiskCount > 0 {
			dialog.ShowConfirm(i18n.T("dialog.high_risk_clean"),
				fmt.Sprintf(i18n.T("dialog.high_risk_msg"), highRiskCount),
				func(confirmed bool) {
					if confirmed {
						a.executeClean(summaryLabel)
					}
				}, a.window)
			return
		}
		a.executeClean(summaryLabel)
	})

	cancelBtn := widget.NewButton(i18n.T("btn.back"), func() {
		a.selectTab(1)
	})

	return container.NewVBox(
		summaryLabel,
		container.NewHBox(confirmBtn, cancelBtn),
	)
}

func (a *App) executeClean(summaryLabel *widget.Label) {
	summaryLabel.SetText(i18n.T("label.cleaning"))
	startTime := time.Now()

	go func() {
		result, err := a.state.RunClean()
		if err != nil {
			fyne.Do(func() {
				dialog.ShowError(err, a.window)
			})
			return
		}

		duration := time.Since(startTime)
		msg := fmt.Sprintf(i18n.T("dialog.clean_result"),
			result.Cleaned, formatSize(result.FreedSize),
			0, formatSize(0),
			result.Failed, duration)
		fyne.Do(func() {
			summaryLabel.SetText(i18n.T("label.clean_complete"))
			dialog.ShowInformation(i18n.T("dialog.clean_done"), msg, a.window)
		})
	}()
}
