package i18n

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

//go:embed embedded/*.json
var embeddedFS embed.FS

var (
	messages = make(map[string]string)
	mu       sync.RWMutex
)

func Init(lang string) {
	mu.Lock()
	defer mu.Unlock()
	load(lang)
}

func load(lang string) {
	// Clear stale keys from previous language
	for k := range messages {
		delete(messages, k)
	}

	// 1. Try embedded locale (always available)
	if data, err := embeddedFS.ReadFile("embedded/" + lang + ".json"); err == nil {
		var m map[string]string
		if err := json.Unmarshal(data, &m); err == nil {
			for k, v := range m {
				messages[k] = v
			}
		}
	}

	// 2. Override with external file if present (allows customization)
	paths := []string{
		"assets/i18n/" + lang + ".json",
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		paths = append(paths, filepath.Join(dir, "assets", "i18n", lang+".json"))
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var m map[string]string
		if err := json.Unmarshal(data, &m); err == nil {
			for k, v := range m {
				messages[k] = v
			}
			return
		}
	}
}

func T(key string) string {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := messages[key]; ok {
		return v
	}
	return key
}

// TDefault returns the translation for key, or fallback if not found.
func TDefault(key, fallback string) string {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := messages[key]; ok {
		return v
	}
	return fallback
}
