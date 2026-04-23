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
)

func (a *App) newConfirmView() fyne.CanvasObject {
	summaryLabel := widget.NewLabel("请先扫描")

	confirmBtn := widget.NewButton("确认清理", func() {
		items := a.state.ScanItems
		if len(items) == 0 {
			dialog.ShowInformation("提示", "没有可清理的项目", a.window)
			return
		}

		// Build clean task
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

		// Check for high-risk items
		var highRiskCount int
		for _, f := range files {
			if f.RiskScore > 60 {
				highRiskCount++
			}
		}
		if highRiskCount > 0 {
			dialog.ShowConfirm("高风险清理",
				fmt.Sprintf("发现 %d 个高风险文件，这些文件将被移入隔离区而非直接删除。是否继续？", highRiskCount),
				func(confirmed bool) {
					if confirmed {
						a.executeClean(files, totalSize, summaryLabel, startTime)
					}
				}, a.window)
			return
		}
		a.executeClean(files, totalSize, summaryLabel, startTime)
	})

	cancelBtn := widget.NewButton("返回", func() {
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

	// Setup executor with quarantine in LOCALAPPDATA
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.TempDir()
	}
	qDir := filepath.Join(localAppData, "CleanMyComputer", "quarantine")
	qm, err := cleaner.NewQuarantineManager(qDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf("创建隔离目录失败: %w", err), a.window)
		return
	}

	// Set OnQuarantine callback for persistence
	qm.OnQuarantine = func(record cleaner.QuarantineRecord) error {
		log.Printf("[Quarantine] %s -> %s (size=%d, expires=%s)",
			record.OriginalPath, record.QuarantinePath, record.Size, record.ExpiresAt.Format(time.DateTime))
		return nil
	}

	executor := cleaner.NewExecutor(qm)

	summaryLabel.SetText("正在清理...")

	go func() {
		result, err := executor.Execute(context.Background(), task)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		duration := time.Since(startTime)
		cleanResult := CleanResult{
			Cleaned:   len(result.Cleaned),
			Failed:    len(result.Failed),
			FreedSize: result.FreedSize,
			Duration:  duration,
		}

		// Save history to database
		a.state.SaveCleanHistory(cleanResult)

		// Notify caller for history persistence
		if a.state.OnCleanComplete != nil {
			a.state.OnCleanComplete(cleanResult)
		}

		msg := fmt.Sprintf("清理完成！\n清理: %d 个文件\n失败: %d 个文件\n释放: %d 字节\n耗时: %v",
			cleanResult.Cleaned, cleanResult.Failed, cleanResult.FreedSize, duration)
		summaryLabel.SetText("清理完成")
		dialog.ShowInformation("清理完成", msg, a.window)

		// Clear scan items
		a.state.ScanItems = nil
	}()
}
