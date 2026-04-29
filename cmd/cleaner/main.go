package main

import (
	"log"
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
	if err := logger.Init(filepath.Join(localAppData, "CleanMyComputer", "logs")); err != nil {
		log.Printf("Warning: failed to init logger: %v", err)
	}

	app := ui.NewApp()
	defer logger.Close()
	app.Run()
}
