package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newRecoveryView() fyne.CanvasObject {
	qList := widget.NewList(
		func() int {
			items, _ := a.state.GetQuarantinedItems()
			return len(items)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil,
				container.NewHBox(widget.NewButton("恢复", nil), widget.NewButton("永久删除", nil)),
				container.NewVBox(widget.NewLabel(""), widget.NewLabel("")),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			items, _ := a.state.GetQuarantinedItems()
			if id >= len(items) {
				return
			}
			item := items[id]
			border := obj.(*fyne.Container)
			vbox := border.Objects[0].(*fyne.Container)
			btnBox := border.Objects[1].(*fyne.Container)

			pathLabel := vbox.Objects[0].(*widget.Label)
			infoLabel := vbox.Objects[1].(*widget.Label)
			pathLabel.SetText(item.OriginalPath)
			infoLabel.SetText(fmt.Sprintf("大小: %s | 隔离时间: %s | 过期: %s",
				formatSize(item.Size),
				item.CreatedAt.Format("2006-01-02 15:04"),
				item.ExpiresAt.Format("2006-01-02 15:04")))

			qPath := item.QuarantinePath
			oPath := item.OriginalPath

			restoreBtn := btnBox.Objects[0].(*widget.Button)
			deleteBtn := btnBox.Objects[1].(*widget.Button)

			restoreBtn.OnTapped = func() {
				dialog.ShowConfirm("恢复文件",
					fmt.Sprintf("将文件恢复到:\n%s", oPath),
					func(ok bool) {
						if !ok {
							return
						}
						if err := a.state.RestoreQuarantinedItem(qPath, oPath); err != nil {
							dialog.ShowError(err, a.window)
							return
						}
						dialog.ShowInformation("成功", "文件已恢复", a.window)
						qList.Refresh()
					}, a.window)
			}

			deleteBtn.OnTapped = func() {
				dialog.ShowConfirm("永久删除",
					"此操作不可逆，文件将被永久删除。确认？",
					func(ok bool) {
						if !ok {
							return
						}
						if err := a.state.DeleteQuarantinedItem(qPath); err != nil {
							dialog.ShowError(err, a.window)
							return
						}
						dialog.ShowInformation("成功", "文件已永久删除", a.window)
						qList.Refresh()
					}, a.window)
			}
		},
	)

	refreshBtn := widget.NewButton("刷新", func() {
		qList.Refresh()
	})

	return container.NewBorder(
		container.NewVBox(widget.NewLabel("恢复中心 - 隔离文件管理"), refreshBtn),
		nil, nil, nil,
		qList,
	)
}
