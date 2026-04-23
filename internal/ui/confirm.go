package ui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

		task := &cleaner.CleanTask{
			Files:     files,
			TotalSize: totalSize,
		}

		// Setup executor with quarantine
		tmpDir := filepath.Join(os.TempDir(), "cleanMyComputer_quarantine")
		qm := cleaner.NewQuarantineManager(tmpDir)
		executor := cleaner.NewExecutor(qm)

		summaryLabel.SetText("正在清理...")

		go func() {
			result, err := executor.Execute(context.Background(), task)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			msg := fmt.Sprintf("清理完成！\n清理: %d 个文件\n失败: %d 个文件\n释放: %d 字节",
				len(result.Cleaned), len(result.Failed), result.FreedSize)
			summaryLabel.SetText("清理完成")
			dialog.ShowInformation("清理完成", msg, a.window)

			// Clear scan items
			a.state.ScanItems = nil
		}()
	})

	cancelBtn := widget.NewButton("返回", func() {
		a.selectTab(1)
	})

	return container.NewVBox(
		summaryLabel,
		container.NewHBox(confirmBtn, cancelBtn),
	)
}
