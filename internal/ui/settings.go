package ui

import (
	"fmt"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

var supportedLangs = []struct {
	Code string
	Name string
}{
	{"zh-CN", "简体中文"},
	{"en-US", "English"},
}

func (a *App) newSettingsView() fyne.CanvasObject {
	var cachedRules []*models.CleanRule
	reloadRules := func() {
		cachedRules = a.state.GetAllRules()
		sort.Slice(cachedRules, func(i, j int) bool {
			return cachedRules[i].RiskScore < cachedRules[j].RiskScore
		})
	}

	var suppressRuleCallback bool

	ruleList := widget.NewList(
		func() int {
			return len(cachedRules)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), widget.NewLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(cachedRules) {
				return
			}
			rule := cachedRules[id]
			hbox := obj.(*fyne.Container)
			check := hbox.Objects[0].(*widget.Check)
			label := hbox.Objects[1].(*widget.Label)

			ruleID := rule.ID
			ruleEnabled := rule.Enabled

			check.OnChanged = func(checked bool) {
				if suppressRuleCallback {
					return
				}
				if err := a.state.SetRuleEnabled(ruleID, checked); err != nil {
					dialog.ShowError(err, a.window)
					return
				}
				rule.Enabled = checked
			}
			suppressRuleCallback = true
			check.SetChecked(ruleEnabled)
			suppressRuleCallback = false

			levelLabel := i18n.T("label.level_safe")
			if rule.Level == 2 {
				levelLabel = i18n.T("label.level_deep")
			} else if rule.Level == 3 {
				levelLabel = i18n.T("label.level_advanced")
			}
			ruleName := i18n.TDefault("rule."+rule.ID+".name", rule.Name)
			ruleDesc := i18n.TDefault("rule."+rule.ID+".desc", rule.Description)
			label.SetText(fmt.Sprintf("%s [%s] %s", ruleName, levelLabel, ruleDesc))
		},
	)

	langSelect := widget.NewSelect(nil, nil)
	var langOptions []string
	var currentLang string
	for _, l := range supportedLangs {
		langOptions = append(langOptions, l.Name)
	}
	langSelect.SetOptions(langOptions)

	if val, err := a.state.GetConfig("language"); err == nil && val != "" {
		currentLang = val
	} else {
		currentLang = "zh-CN"
	}
	for i, l := range supportedLangs {
		if l.Code == currentLang {
			langSelect.SetSelectedIndex(i)
			break
		}
	}

	langLabel := widget.NewLabel(i18n.T("label.language"))

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
		idx := langSelect.SelectedIndex()
		langChanged := false
		if idx >= 0 && idx < len(supportedLangs) {
			langCode := supportedLangs[idx].Code
			oldLang, _ := a.state.GetConfig("language")
			if oldLang != langCode {
				langChanged = true
			}
			if err := a.state.SetConfig("language", langCode); err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			i18n.Init(langCode)
		}

		autoEnabled := "false"
		if autoCleanCheck.Checked {
			autoEnabled = "true"
			if intervalEntry.Text != "" {
				if h, err := strconv.Atoi(intervalEntry.Text); err != nil || h <= 0 {
					dialog.ShowError(fmt.Errorf("invalid interval: %s", intervalEntry.Text), a.window)
					return
				}
			}
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

		if langChanged {
			a.buildTabs()
			return
		}

		dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("dialog.settings_saved"), a.window)
	})

	refreshBtn := widget.NewButton(i18n.T("btn.refresh_rules"), func() {
		reloadRules()
		ruleList.Refresh()
	})

	ruleSelectAllBtn := widget.NewButton(i18n.T("btn.select_all"), func() {
		for _, r := range cachedRules {
			a.state.SetRuleEnabled(r.ID, true)
		}
		reloadRules()
		ruleList.Refresh()
	})

	ruleDeselectAllBtn := widget.NewButton(i18n.T("btn.deselect_all"), func() {
		for _, r := range cachedRules {
			a.state.SetRuleEnabled(r.ID, false)
		}
		reloadRules()
		ruleList.Refresh()
	})

	reloadRules()

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(i18n.T("label.settings"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewHBox(langLabel, langSelect),
			widget.NewSeparator(),
			container.NewHBox(autoCleanCheck, intervalLabel, intervalEntry),
			saveBtn,
			widget.NewSeparator(),
			widget.NewLabel(i18n.T("label.rule_toggle")),
			container.NewHBox(refreshBtn, ruleSelectAllBtn, ruleDeselectAllBtn),
		),
		nil, nil, nil,
		ruleList,
	)
}

func (a *App) applyAutoCleanSettings(enabled bool, intervalStr string) {
	if a.scheduler != nil {
		a.scheduler.Stop()
		a.scheduler = nil
	}

	if !enabled {
		return
	}

	hours := 24
	if intervalStr != "" {
		if h, err := strconv.Atoi(intervalStr); err == nil && h > 0 {
			hours = h
		}
	}

	a.scheduler = newScheduler(a.state.Orchestrator)
	a.scheduler.SetInterval(hours)
	a.scheduler.Start()
}
