package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

type rowKind int

const (
	rowGroup rowKind = iota
	rowItem
)

type displayRow struct {
	kind      rowKind
	ruleID    string
	ruleName  string
	count     int
	totalSize int64
	item      *models.ScanItem
}

type tappableLabel struct {
	widget.Label
	onTapped func()
}

func newTappableLabel(text string) *tappableLabel {
	l := &tappableLabel{}
	l.ExtendBaseWidget(l)
	l.SetText(text)
	return l
}

func (t *tappableLabel) Tapped(*fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func (t *tappableLabel) TappedSecondary(*fyne.PointEvent) {}

type scannerPage struct {
	app              *App
	rows             []displayRow
	allItems         []*models.ScanItem
	checked          map[string]bool
	collapsed        map[string]bool
	filterRisk       string
	minSizeBytes     int64
	resultList       *widget.List
	summaryLabel     *widget.Label
	suppressCallback bool
}

func newScannerPage(app *App) *scannerPage {
	return &scannerPage{
		app:        app,
		checked:    make(map[string]bool),
		collapsed:  make(map[string]bool),
		filterRisk: "all",
	}
}

func (p *scannerPage) buildUI() fyne.CanvasObject {
	headerLabel := widget.NewLabelWithStyle(
		i18n.T("label.scan_results"),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	riskOptions := []string{
		i18n.T("filter.all"),
		i18n.T("risk.safe"),
		i18n.T("risk.moderate"),
		i18n.T("risk.high"),
	}
	riskValues := []string{"all", "safe", "moderate", "high"}
	riskSelect := widget.NewSelect(riskOptions, func(selected string) {
		for i, opt := range riskOptions {
			if opt == selected {
				p.filterRisk = riskValues[i]
				break
			}
		}
		p.rebuildRows()
		p.resultList.Refresh()
		p.updateSummary()
	})

	sizeLabel := widget.NewLabel(i18n.T("label.filter_size"))
	sizeEntry := widget.NewEntry()
	sizeEntry.SetPlaceHolder("0")
	sizeEntry.Validator = nil
	sizeEntry.OnChanged = func(val string) {
		kb := 0
		if val != "" {
			if n, err := strconv.Atoi(val); err == nil && n >= 0 {
				kb = n
			}
		}
		newMin := int64(kb) * 1024
		if newMin != p.minSizeBytes {
			p.minSizeBytes = newMin
			p.rebuildRows()
			p.resultList.Refresh()
			p.updateSummary()
		}
	}

	selectAllBtn := widget.NewButton(i18n.T("btn.select_all"), func() {
		for _, item := range p.allItems {
			if p.matchesFilter(item) {
				p.checked[item.Path] = true
			}
		}
		p.resultList.Refresh()
		p.updateSummary()
	})

	deselectAllBtn := widget.NewButton(i18n.T("btn.deselect_all"), func() {
		p.checked = make(map[string]bool)
		p.resultList.Refresh()
		p.updateSummary()
	})

	p.summaryLabel = widget.NewLabel("")

	p.resultList = widget.NewList(
		func() int {
			return len(p.rows)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), newTappableLabel(""))
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(p.rows) {
				return
			}
			row := p.rows[id]
			hbox := obj.(*fyne.Container)
			check := hbox.Objects[0].(*widget.Check)
			lbl := hbox.Objects[1].(*tappableLabel)

			if row.kind == rowGroup {
				indicator := "▾"
				if p.collapsed[row.ruleID] {
					indicator = "▸"
				}
				lbl.SetText(fmt.Sprintf("%s %s (%d %s, %s)",
					indicator, row.ruleName, row.count, i18n.T("label.file_count"), formatSize(row.totalSize)))
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				lbl.Refresh()

				capturedRuleID := row.ruleID
				lbl.onTapped = func() {
					p.collapsed[capturedRuleID] = !p.collapsed[capturedRuleID]
					p.rebuildRows()
					p.resultList.Refresh()
					p.updateSummary()
				}
			} else {
				item := row.item
				riskStr := riskScoreToLabel(item.RiskScore)
				displayPath := item.Path
				if item.Type == "command" {
					displayPath = i18n.T("label.command") + " " + item.Path
				}
				lbl.SetText(fmt.Sprintf("    %s | %s: %s | %s: %s (%d)",
					displayPath,
					i18n.T("label.size"), formatSize(item.Size),
					i18n.T("label.risk"), riskStr, item.RiskScore))
				lbl.TextStyle = fyne.TextStyle{}
				lbl.Refresh()

				lbl.onTapped = nil
			}

			isChecked := false
			if row.kind == rowGroup {
				isChecked = p.isGroupAllChecked(row.ruleID)
			} else if row.item != nil {
				isChecked = p.checked[row.item.Path]
			}

			capturedKind := row.kind
			capturedRuleID := row.ruleID
			var capturedPath string
			if row.item != nil {
				capturedPath = row.item.Path
			}

			check.OnChanged = func(c bool) {
				if p.suppressCallback {
					return
				}
				if capturedKind == rowGroup {
					p.setGroupChecked(capturedRuleID, c)
				} else {
					p.checked[capturedPath] = c
				}
				p.resultList.Refresh()
				p.updateSummary()
			}

			p.suppressCallback = true
			check.SetChecked(isChecked)
			p.suppressCallback = false
		},
	)
	riskSelect.SetSelectedIndex(0)

	cleanBtn := widget.NewButton(i18n.T("btn.start_clean"), func() {
		p.startClean()
	})

	return container.NewBorder(
		container.NewVBox(
			headerLabel,
			container.NewHBox(
				widget.NewLabel(i18n.T("label.filter_risk")),
				riskSelect,
				sizeLabel,
				sizeEntry,
			),
			container.NewHBox(selectAllBtn, deselectAllBtn, p.summaryLabel),
		),
		cleanBtn,
		nil, nil,
		p.resultList,
	)
}

func (p *scannerPage) refreshData() {
	p.allItems = p.app.state.GetScanItemsSafe()
	p.checked = make(map[string]bool)
	p.collapsed = make(map[string]bool)
	p.rebuildRows()
	if p.resultList != nil {
		p.resultList.Refresh()
	}
	p.updateSummary()
}

func (p *scannerPage) rebuildRows() {
	type group struct {
		ruleID    string
		ruleName  string
		items     []*models.ScanItem
		totalSize int64
	}
	groups := make(map[string]*group)
	var groupOrder []string

	for _, item := range p.allItems {
		if !p.matchesFilter(item) {
			continue
		}
		ruleID := item.RuleID
		if ruleID == "" {
			ruleID = "unknown"
		}
		if _, ok := groups[ruleID]; !ok {
			ruleName := i18n.TDefault("rule."+ruleID+".name", ruleID)
			groups[ruleID] = &group{ruleID: ruleID, ruleName: ruleName}
			groupOrder = append(groupOrder, ruleID)
		}
		g := groups[ruleID]
		g.items = append(g.items, item)
		g.totalSize += item.Size
	}

	p.rows = nil
	for _, ruleID := range groupOrder {
		g := groups[ruleID]
		p.rows = append(p.rows, displayRow{
			kind:      rowGroup,
			ruleID:    g.ruleID,
			ruleName:  g.ruleName,
			count:     len(g.items),
			totalSize: g.totalSize,
		})
		if !p.collapsed[ruleID] {
			for _, item := range g.items {
				p.rows = append(p.rows, displayRow{
					kind:   rowItem,
					ruleID: g.ruleID,
					item:   item,
				})
			}
		}
	}
}

func (p *scannerPage) matchesFilter(item *models.ScanItem) bool {
	switch p.filterRisk {
	case "safe":
		if item.RiskScore > 30 {
			return false
		}
	case "moderate":
		if item.RiskScore <= 30 || item.RiskScore > 60 {
			return false
		}
	case "high":
		if item.RiskScore <= 60 {
			return false
		}
	}
	if p.minSizeBytes > 0 && item.Size < p.minSizeBytes {
		return false
	}
	return true
}

func (p *scannerPage) isGroupAllChecked(ruleID string) bool {
	count := 0
	checked := 0
	for _, item := range p.allItems {
		if item.RuleID == ruleID && p.matchesFilter(item) {
			count++
			if p.checked[item.Path] {
				checked++
			}
		}
	}
	return count > 0 && count == checked
}

func (p *scannerPage) setGroupChecked(ruleID string, val bool) {
	for _, item := range p.allItems {
		if item.RuleID == ruleID && p.matchesFilter(item) {
			p.checked[item.Path] = val
		}
	}
}

func (p *scannerPage) getCheckedItems() []*models.ScanItem {
	var items []*models.ScanItem
	for _, item := range p.allItems {
		if p.checked[item.Path] {
			items = append(items, item)
		}
	}
	return items
}

func (p *scannerPage) updateSummary() {
	items := p.getCheckedItems()
	var totalSize int64
	for _, item := range items {
		totalSize += item.Size
	}
	p.summaryLabel.SetText(fmt.Sprintf(i18n.T("label.selected_info"), len(items), formatSize(totalSize)))
}

func (p *scannerPage) startClean() {
	items := p.getCheckedItems()
	if len(items) == 0 {
		dialog.ShowInformation(i18n.T("dialog.tip"), i18n.T("label.no_items"), p.app.window)
		return
	}

	var highRiskCount int
	for _, item := range items {
		if item.RiskScore > 60 {
			highRiskCount++
		}
	}

	doClean := func() {
		p.app.state.SetScanItemsForClean(items)

		p.summaryLabel.SetText(i18n.T("label.cleaning"))
		startTime := time.Now()

		go func() {
			result, err := p.app.state.RunClean()
			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(err, p.app.window)
				})
				return
			}

			duration := time.Since(startTime)
			msg := fmt.Sprintf(i18n.T("dialog.clean_result"),
				result.Cleaned, formatSize(result.FreedSize),
				result.Failed, duration)
			fyne.Do(func() {
				p.summaryLabel.SetText(i18n.T("label.clean_complete"))
				dialog.ShowInformation(i18n.T("dialog.clean_done"), msg, p.app.window)
				p.refreshData()
			})
		}()
	}

	if highRiskCount > 0 {
		dialog.ShowConfirm(i18n.T("dialog.high_risk_clean"),
			fmt.Sprintf(i18n.T("dialog.high_risk_msg"), highRiskCount),
			func(confirmed bool) {
				if confirmed {
					doClean()
				}
			}, p.app.window)
		return
	}
	doClean()
}

func riskScoreToLabel(score int) string {
	if score > 60 {
		return i18n.T("risk.high")
	}
	if score > 30 {
		return i18n.T("risk.moderate")
	}
	return i18n.T("risk.safe")
}
