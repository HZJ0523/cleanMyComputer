package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/core/report"
	"github.com/hzj0523/cleanMyComputer/internal/models"
)

func (a *App) newHistoryView() fyne.CanvasObject {
	historyList := widget.NewList(
		func() int {
			records, _ := a.state.GetHistory()
			if records == nil {
				return 0
			}
			return len(records)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel(""), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			records, _ := a.state.GetHistory()
			if id >= len(records) {
				return
			}
			r := records[id]
			vbox := obj.(*fyne.Container)
			line1 := vbox.Objects[0].(*widget.Label)
			line2 := vbox.Objects[1].(*widget.Label)
			line1.SetText(fmt.Sprintf("%s | 释放: %s | 状态: %s",
				r.StartTime.Format("2006-01-02 15:04"), formatSize(r.FreedSize), r.Status))
			line2.SetText(fmt.Sprintf("文件数: %d | 耗时: %s",
				r.TotalFiles, r.EndTime.Sub(r.StartTime).Round(time.Second)))
		},
	)

	refreshBtn := widget.NewButton("刷新", func() {
		historyList.Refresh()
	})

	exportBtn := widget.NewButton("导出报告", func() {
		records, err := a.state.GetHistory()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if len(records) == 0 {
			dialog.ShowInformation("提示", "没有历史记录可导出", a.window)
			return
		}

		gen := report.NewGenerator()

		var items []*models.ScanItem
		var totalFreed int64
		for _, rec := range records {
			totalFreed += rec.FreedSize
			items = append(items, &models.ScanItem{
				Path:    fmt.Sprintf("清理记录 #%d", rec.ID),
				Size:    rec.FreedSize,
				ModTime: rec.StartTime,
			})
		}

		r := &report.Report{
			ScanLevel: records[0].ScanLevel,
			StartTime: records[0].StartTime,
			EndTime:   records[len(records)-1].EndTime,
			FreedSize: totalFreed,
			Items:     items,
		}

		home, _ := os.UserHomeDir()
		exportPath := filepath.Join(home, "Desktop", fmt.Sprintf("清理报告_%s.txt", time.Now().Format("20060102_150405")))
		if err := gen.ExportToFile(r, exportPath); err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		dialog.ShowInformation("导出成功", fmt.Sprintf("报告已保存到:\n%s", exportPath), a.window)
	})

	return container.NewBorder(
		container.NewVBox(widget.NewLabel("清理历史"), container.NewHBox(refreshBtn, exportBtn)),
		nil, nil, nil,
		historyList,
	)
}
