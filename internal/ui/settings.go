package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) newSettingsView() fyne.CanvasObject {
	// Rule list
	ruleList := widget.NewList(
		func() int {
			rules := a.state.GetAllRules()
			return len(rules)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			rules := a.state.GetAllRules()
			if id >= len(rules) {
				return
			}
			rule := rules[id]
			hbox := obj.(*fyne.Container)
			check := hbox.Objects[0].(*widget.Check)
			label := hbox.Objects[1].(*widget.Label)

			check.SetChecked(rule.Enabled)
			check.OnChanged = func(checked bool) {
				if err := a.state.SetRuleEnabled(rule.ID, checked); err != nil {
					dialog.ShowError(err, a.window)
				}
			}
			levelLabel := "安全"
			if rule.Level == 2 {
				levelLabel = "深度"
			} else if rule.Level == 3 {
				levelLabel = "高级"
			}
			label.SetText(fmt.Sprintf("%s [%s] %s", rule.Name, levelLabel, rule.Description))
		},
	)

	// Retention hours
	retentionEntry := widget.NewEntry()
	retentionEntry.SetPlaceHolder("24")
	if val, err := a.state.GetConfig("quarantine_retention_hours"); err == nil && val != "" {
		retentionEntry.SetText(val)
	}

	retentionLabel := widget.NewLabel("隔离区保留时间（小时）")
	saveBtn := widget.NewButton("保存设置", func() {
		if retentionEntry.Text != "" {
			if err := a.state.SetConfig("quarantine_retention_hours", retentionEntry.Text); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
		}
		dialog.ShowInformation("提示", "设置已保存", a.window)
	})

	refreshBtn := widget.NewButton("刷新规则列表", func() {
		ruleList.Refresh()
	})

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("设置", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewHBox(retentionLabel, retentionEntry),
			saveBtn,
			widget.NewSeparator(),
			widget.NewLabel("清理规则启用/禁用"),
			refreshBtn,
		),
		nil, nil, nil,
		ruleList,
	)
}
