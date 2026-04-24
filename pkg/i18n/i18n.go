package i18n

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	messages = make(map[string]string)
	once     sync.Once
)

func Init(lang string) {
	once.Do(func() {
		load(lang)
	})
}

func load(lang string) {
	paths := []string{
		"assets/i18n/" + lang + ".json",
	}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, exe[:len(exe)-len("/cleanMyComputer.exe")]+"assets/i18n/"+lang+".json")
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
	if v, ok := messages[key]; ok {
		return v
	}
	return key
}
