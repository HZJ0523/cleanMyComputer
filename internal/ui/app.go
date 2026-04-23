package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type App struct {
	fyneApp fyne.App
	window  fyne.Window
}

func NewApp() *App {
	a := app.New()
	w := a.NewWindow("CleanMyComputer")
	w.Resize(fyne.NewSize(1024, 768))

	return &App{fyneApp: a, window: w}
}

func (a *App) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem("首页", NewDashboard()),
		container.NewTabItem("扫描结果", NewScannerView()),
		container.NewTabItem("清理确认", NewConfirmView()),
		container.NewTabItem("历史记录", NewHistoryView()),
		container.NewTabItem("设置", NewSettingsView()),
	)
	a.window.SetContent(tabs)
	a.window.ShowAndRun()
}
