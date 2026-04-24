package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	return container.NewBorder(
		container.NewVBox(widget.NewLabel("清理历史"), refreshBtn),
		nil, nil, nil,
		historyList,
	)
}
