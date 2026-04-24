package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

func (a *App) newSettingsView() fyne.CanvasObject {
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
			levelLabel := i18n.T("label.level_safe")
			if rule.Level == 2 {
				levelLabel = i18n.T("label.level_deep")
			} else if rule.Level == 3 {
				levelLabel = i18n.T("label.level_advanced")
			}
			label.SetText(fmt.Sprintf("%s [%s] %s", rule.Name, levelLabel, rule.Description))
		},
	)

	retentionEntry := widget.NewEntry()
	retentionEntry.SetPlaceHolder("24")
	if val, err := a.state.GetConfig("quarantine_retention_hours"); err == nil && val != "" {
		retentionEntry.SetText(val)
	}
	retentionLabel := widget.NewLabel(i18n.T("label.quarantine_retention"))

	autoCleanCheck := widget.NewCheck(i18n.T("label.auto_clean_enable"), nil)
	if val, err := a.state.GetConfig("auto_clean_enabled"); err == nil && val == "true" {
		autoCleanCheck.SetChecked(true)
	}

	intervalEntry := widget.NewEntry()
	intervalEntry.SetPlaceHolder("24")
	if val, err := a.state.GetConfig("auto_clean_interval_hours"); err == nil && val != "" {
		intervalEntry.SetText(val)
	}
	intervalLabel := widget.NewLabel(i18n.T("label.auto_clean_interval"))

	saveBtn := widget.NewButton(i18n.T("btn.save"), func() {
		if retentionEntry.Text != "" {
			if err := a.state.SetConfig("quarantine_retention_hours", retentionEntry.Text); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
		}

		autoEnabled := "false"
		if autoCleanCheck.Checked {
			autoEnabled = "true"
		}
		if err := a.state.SetConfig("auto_clean_enabled", autoEnabled); err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		if intervalEntry.Text != "" {
			if err := a.state.SetConfig("auto_clean_interval_hours", intervalEntry.Text); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
		}

		a.applyAutoCleanSettings(autoCleanCheck.Checked, intervalEntry.Text)

		dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("dialog.settings_saved"), a.window)
	})

	refreshBtn := widget.NewButton(i18n.T("btn.refresh_rules"), func() {
		ruleList.Refresh()
	})

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(i18n.T("label.settings"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewHBox(retentionLabel, retentionEntry),
			widget.NewSeparator(),
			container.NewHBox(autoCleanCheck, intervalLabel, intervalEntry),
			saveBtn,
			widget.NewSeparator(),
			widget.NewLabel(i18n.T("label.rule_toggle")),
			refreshBtn,
		),
		nil, nil, nil,
		ruleList,
	)
}

func (a *App) applyAutoCleanSettings(enabled bool, intervalStr string) {
	if !enabled {
		if a.scheduler != nil {
			a.scheduler.Stop()
		}
		return
	}

	hours := 24
	if intervalStr != "" {
		if h, err := strconv.Atoi(intervalStr); err == nil && h > 0 {
			hours = h
		}
	}

	if a.scheduler == nil {
		a.scheduler = newScheduler(a.state.Orchestrator)
	}
	a.scheduler.SetInterval(hours)
	a.scheduler.Start()
}
