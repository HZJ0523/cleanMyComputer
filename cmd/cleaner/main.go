package main

import (
	"os"
	"path/filepath"

	"github.com/hzj0523/cleanMyComputer/internal/ui"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
	"github.com/hzj0523/cleanMyComputer/pkg/logger"
)

func main() {
	i18n.Init("zh-CN")

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.TempDir()
	}
	logger.Init(filepath.Join(localAppData, "CleanMyComputer", "logs"))

	app := ui.NewApp()
	app.Run()
}
