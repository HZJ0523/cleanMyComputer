package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newScannerView() fyne.CanvasObject {
	headerLabel := widget.NewLabel("扫描结果")

	resultList := widget.NewList(
		func() int {
			return len(a.state.ScanItems)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel(""), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			items := a.state.ScanItems
			if id < len(items) {
				item := items[id]
				vbox := obj.(*fyne.Container)
				pathLabel := vbox.Objects[0].(*widget.Label)
				infoLabel := vbox.Objects[1].(*widget.Label)
				pathLabel.SetText(item.Path)
				riskLabel := "安全"
				if item.RiskScore > 60 {
					riskLabel = "高风险"
				} else if item.RiskScore > 30 {
					riskLabel = "中等风险"
				}
				infoLabel.SetText(fmt.Sprintf("大小: %d 字节 | 风险: %s (%d)", item.Size, riskLabel, item.RiskScore))
			}
		},
	)

	cleanBtn := widget.NewButton("开始清理", func() {
		a.selectTab(2)
	})

	return container.NewBorder(
		headerLabel,
		cleanBtn,
		nil, nil,
		resultList,
	)
}
