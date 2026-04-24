package ui

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
)

type App struct {
	fyneApp   fyne.App
	window    fyne.Window
	state     *AppState
	tabs      *container.AppTabs
	scheduler *scheduler
}

func NewApp() *App {
	a := app.New()
	w := a.NewWindow("CleanMyComputer")
	w.Resize(fyne.NewSize(1024, 768))

	state := NewAppState()

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.TempDir()
	}
	dbPath := filepath.Join(localAppData, "CleanMyComputer", "cleaner.db")
	if err := state.InitDB(dbPath); err != nil {
		log.Printf("Warning: failed to init database: %v", err)
	}

	return &App{fyneApp: a, window: w, state: state}
}

func (a *App) Run() {
	defer a.state.CloseDB()
	defer func() {
		if a.scheduler != nil {
			a.scheduler.Stop()
		}
	}()

	a.state.CleanupExpiredQuarantine()
	a.restoreScheduler()

	dashboard := a.newDashboard()
	scannerView := a.newScannerView()
	confirmView := a.newConfirmView()
	recoveryView := a.newRecoveryView()
	historyView := a.newHistoryView()
	settingsView := a.newSettingsView()

	a.tabs = container.NewAppTabs(
		container.NewTabItem(i18n.T("tab.dashboard"), dashboard),
		container.NewTabItem(i18n.T("tab.scan"), scannerView),
		container.NewTabItem(i18n.T("tab.confirm"), confirmView),
		container.NewTabItem(i18n.T("tab.recovery"), recoveryView),
		container.NewTabItem(i18n.T("tab.history"), historyView),
		container.NewTabItem(i18n.T("tab.settings"), settingsView),
	)
	a.window.SetContent(a.tabs)
	a.window.ShowAndRun()
}

func (a *App) selectTab(index int) {
	if a.tabs != nil {
		fyne.Do(func() {
			a.tabs.SelectIndex(index)
		})
	}
}

func (a *App) restoreScheduler() {
	enabled, err := a.state.GetConfig("auto_clean_enabled")
	if err != nil || enabled != "true" {
		return
	}
	intervalStr, _ := a.state.GetConfig("auto_clean_interval_hours")
	hours := 24
	if intervalStr != "" {
		if h, err := strconv.Atoi(intervalStr); err == nil && h > 0 {
			hours = h
		}
	}
	a.scheduler = newScheduler(a.state.Orchestrator)
	a.scheduler.SetInterval(hours)
	a.scheduler.Start()
	log.Printf("[App] auto-clean scheduler restored (interval=%dh)", hours)
}
