package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/core/cleaner"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func (a *App) newConfirmView() fyne.CanvasObject {
	summaryLabel := widget.NewLabel(i18n.T("label.please_scan"))

	confirmBtn := widget.NewButton(i18n.T("btn.confirm_clean"), func() {
		items := a.state.ScanItems
		if len(items) == 0 {
			dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("label.no_items"), a.window)
			return
		}

		var files []*cleaner.FileItem
		var totalSize int64
		for _, item := range items {
			files = append(files, &cleaner.FileItem{
				Path:      item.Path,
				Size:      item.Size,
				RiskScore: item.RiskScore,
			})
			totalSize += item.Size
		}

		startTime := time.Now()

		var highRiskCount int
		for _, f := range files {
			if f.RiskScore > 60 {
				highRiskCount++
			}
		}
		if highRiskCount > 0 {
			dialog.ShowConfirm(i18n.T("dialog.high_risk_clean"),
				fmt.Sprintf(i18n.T("dialog.high_risk_msg"), highRiskCount),
				func(confirmed bool) {
					if confirmed {
						a.executeClean(files, totalSize, summaryLabel, startTime)
					}
				}, a.window)
			return
		}
		a.executeClean(files, totalSize, summaryLabel, startTime)
	})

	cancelBtn := widget.NewButton(i18n.T("btn.back"), func() {
		a.selectTab(1)
	})

	return container.NewVBox(
		summaryLabel,
		container.NewHBox(confirmBtn, cancelBtn),
	)
}

func (a *App) executeClean(files []*cleaner.FileItem, totalSize int64, summaryLabel *widget.Label, startTime time.Time) {
	task := &cleaner.CleanTask{
		Files:     files,
		TotalSize: totalSize,
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.TempDir()
	}
	qDir := filepath.Join(localAppData, "CleanMyComputer", "quarantine")
	qm, err := cleaner.NewQuarantineManager(qDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf("%s: %w", i18n.T("dialog.create_quarantine_failed"), err), a.window)
		return
	}

	qm.OnQuarantine = func(record cleaner.QuarantineRecord) error {
		log.Printf("[Quarantine] %s -> %s (size=%d, expires=%s)",
			record.OriginalPath, record.QuarantinePath, record.Size, record.ExpiresAt.Format(time.DateTime))
		if err := a.state.SaveQuarantineRecord(record); err != nil {
			log.Printf("Failed to save quarantine record: %v", err)
		}
		return nil
	}

	executor := cleaner.NewExecutor(qm)

	summaryLabel.SetText(i18n.T("label.cleaning"))

	go func() {
		result, err := executor.Execute(context.Background(), task)
		if err != nil {
			fyne.Do(func() {
				dialog.ShowError(err, a.window)
			})
			return
		}

		duration := time.Since(startTime)
		cleanResult := CleanResult{
			Cleaned:   len(result.Cleaned),
			Failed:    len(result.Failed),
			FreedSize: result.FreedSize,
			Duration:  duration,
		}

		a.state.SaveCleanHistory(cleanResult)

		msg := fmt.Sprintf(i18n.T("dialog.clean_result"),
			len(result.Cleaned)-len(result.Quarantined), formatSize(result.FreedSize),
			len(result.Quarantined), formatSize(result.QuarantinedSize),
			len(result.Failed), duration)
		fyne.Do(func() {
			summaryLabel.SetText(i18n.T("label.clean_complete"))
			dialog.ShowInformation(i18n.T("dialog.clean_done"), msg, a.window)
			a.state.ScanItems = nil
		})
	}()
}
