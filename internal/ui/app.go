package ui

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type App struct {
	fyneApp  fyne.App
	window   fyne.Window
	state    *AppState
	tabs     *container.AppTabs
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

	dashboard := a.newDashboard()
	scannerView := a.newScannerView()
	confirmView := a.newConfirmView()
	recoveryView := a.newRecoveryView()
	historyView := a.newHistoryView()
	settingsView := a.newSettingsView()

	a.tabs = container.NewAppTabs(
		container.NewTabItem("首页", dashboard),
		container.NewTabItem("扫描结果", scannerView),
		container.NewTabItem("清理确认", confirmView),
		container.NewTabItem("恢复中心", recoveryView),
		container.NewTabItem("历史记录", historyView),
		container.NewTabItem("设置", settingsView),
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
