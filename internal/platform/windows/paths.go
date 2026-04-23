package windows

import (
	"os"
	"path/filepath"
)

type PathDiscovery struct{}

func NewPathDiscovery() *PathDiscovery {
	return &PathDiscovery{}
}

func (p *PathDiscovery) BrowserCachePaths() []string {
	localAppData := os.Getenv("LOCALAPPDATA")

	return []string{
		// Chrome
		filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "Cache"),
		// Edge
		filepath.Join(localAppData, "Microsoft", "Edge", "User Data", "Default", "Cache"),
		// Firefox
		filepath.Join(localAppData, "Mozilla", "Firefox", "Profiles"),
	}
}

func (p *PathDiscovery) SystemTempPaths() []string {
	return []string{
		os.Getenv("TEMP"),
		os.Getenv("TMP"),
		filepath.Join(os.Getenv("SystemRoot"), "Temp"),
	}
}

func (p *PathDiscovery) DevCachePaths() map[string]string {
	localAppData := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")
	home, _ := os.UserHomeDir()

	return map[string]string{
		"npm":    filepath.Join(appData, "npm-cache"),
		"pip":    filepath.Join(localAppData, "pip", "Cache"),
		"maven":  filepath.Join(home, ".m2", "repository"),
		"gradle": filepath.Join(home, ".gradle", "caches"),
		"cargo":  filepath.Join(home, ".cargo", "registry"),
		"go":     filepath.Join(os.Getenv("GOPATH"), "pkg", "mod"),
	}
}
