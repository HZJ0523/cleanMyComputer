package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

type recoveryView struct {
	app  *App
	list *widget.List
}

func (r *recoveryView) build() fyne.CanvasObject {
	r.list = widget.NewList(
		func() int {
			items, _ := r.app.state.GetQuarantinedItems()
			return len(items)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil,
				container.NewHBox(widget.NewButton(i18n.T("btn.restore"), nil), widget.NewButton(i18n.T("btn.permanent_delete"), nil)),
				container.NewVBox(widget.NewLabel(""), widget.NewLabel("")),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			items, _ := r.app.state.GetQuarantinedItems()
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
			infoLabel.SetText(fmt.Sprintf("%s: %s | %s: %s | %s: %s",
				i18n.T("label.size"), formatSize(item.Size),
				i18n.T("label.quarantine_time"), item.CreatedAt.Format("2006-01-02 15:04"),
				i18n.T("label.expires"), item.ExpiresAt.Format("2006-01-02 15:04")))

			qPath := item.QuarantinePath
			oPath := item.OriginalPath

			restoreBtn := btnBox.Objects[0].(*widget.Button)
			deleteBtn := btnBox.Objects[1].(*widget.Button)

			restoreBtn.OnTapped = func() {
				dialog.ShowConfirm(i18n.T("dialog.confirm_restore"),
					fmt.Sprintf(i18n.T("dialog.restore_msg"), oPath),
					func(ok bool) {
						if !ok {
							return
						}
						if err := r.app.state.RestoreQuarantinedItem(qPath, oPath); err != nil {
							dialog.ShowError(err, r.app.window)
							return
						}
						dialog.ShowInformation(i18n.T("dialog.success"), i18n.T("dialog.file_restored"), r.app.window)
						r.list.Refresh()
					}, r.app.window)
			}

			deleteBtn.OnTapped = func() {
				dialog.ShowConfirm(i18n.T("dialog.confirm_permanent_delete"),
					i18n.T("dialog.permanent_delete_msg"),
					func(ok bool) {
						if !ok {
							return
						}
						if err := r.app.state.DeleteQuarantinedItem(qPath); err != nil {
							dialog.ShowError(err, r.app.window)
							return
						}
						dialog.ShowInformation(i18n.T("dialog.success"), i18n.T("dialog.file_deleted"), r.app.window)
						r.list.Refresh()
					}, r.app.window)
			}
		},
	)

	refreshBtn := widget.NewButton(i18n.T("btn.refresh"), func() {
		r.list.Refresh()
	})

	return container.NewBorder(
		container.NewVBox(widget.NewLabel(i18n.T("label.recovery_center")), refreshBtn),
		nil, nil, nil,
		r.list,
	)
}

func (a *App) newRecoveryView() fyne.CanvasObject {
	v := &recoveryView{app: a}
	return v.build()
}
