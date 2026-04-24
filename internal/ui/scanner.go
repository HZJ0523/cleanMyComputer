package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func (a *App) newScannerView() fyne.CanvasObject {
	headerLabel := widget.NewLabel(i18n.T("label.scan_results"))

	resultList := widget.NewList(
		func() int {
			return a.state.GetScanItemCount()
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel(""), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			items := a.state.GetScanItemsSafe()
			if id < len(items) {
				item := items[id]
				vbox := obj.(*fyne.Container)
				pathLabel := vbox.Objects[0].(*widget.Label)
				infoLabel := vbox.Objects[1].(*widget.Label)

				displayPath := item.Path
				if item.Size == 0 && item.ModTime.IsZero() {
					displayPath = i18n.T("label.command") + " " + item.Path
				}
				pathLabel.SetText(displayPath)

				riskLabel := i18n.T("risk.safe")
				if item.RiskScore > 60 {
					riskLabel = i18n.T("risk.high")
				} else if item.RiskScore > 30 {
					riskLabel = i18n.T("risk.moderate")
				}
				infoLabel.SetText(fmt.Sprintf("%s: %s | %s: %s (%d)",
					i18n.T("label.size"), formatSize(item.Size),
					i18n.T("label.risk"), riskLabel, item.RiskScore))
			}
		},
	)

	cleanBtn := widget.NewButton(i18n.T("btn.start_clean"), func() {
		a.selectTab(2)
	})

	return container.NewBorder(
		headerLabel,
		cleanBtn,
		nil, nil,
		resultList,
	)
}
