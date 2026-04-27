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
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
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
			line1.SetText(fmt.Sprintf("%s | %s: %s | %s: %s",
				r.StartTime.Format("2006-01-02 15:04"),
				i18n.T("label.released"), formatSize(r.FreedSize),
				i18n.T("label.status"), r.Status))
			line2.SetText(fmt.Sprintf("%s: %d | %s: %s",
				i18n.T("label.file_count"), r.TotalFiles,
				i18n.T("label.duration"), r.EndTime.Sub(r.StartTime).Round(time.Second)))
		},
	)

	refreshBtn := widget.NewButton(i18n.T("btn.refresh"), func() {
		historyList.Refresh()
	})

	exportBtn := widget.NewButton(i18n.T("btn.export_report"), func() {
		records, err := a.state.GetHistory()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if len(records) == 0 {
			dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("dialog.no_history"), a.window)
			return
		}

		gen := report.NewGenerator()

		var items []*models.ScanItem
		var totalFreed int64
		for _, rec := range records {
			totalFreed += rec.FreedSize
			items = append(items, &models.ScanItem{
				Path:    fmt.Sprintf(i18n.T("label.clean_record"), rec.ID),
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
		exportPath := filepath.Join(home, fmt.Sprintf(i18n.T("report.filename"), time.Now().Format("20060102_150405")))
		if err := gen.ExportToFile(r, exportPath); err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		dialog.ShowInformation(i18n.T("dialog.export_success"),
			fmt.Sprintf(i18n.T("dialog.export_msg"), exportPath), a.window)
	})

	return container.NewBorder(
		container.NewVBox(widget.NewLabel(i18n.T("label.clean_history")), container.NewHBox(refreshBtn, exportBtn)),
		nil, nil, nil,
		historyList,
	)
}
