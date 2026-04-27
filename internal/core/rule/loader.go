package rule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Loader struct {
	rulesDir string
}

func NewLoader() *Loader {
	exe, err := os.Executable()
	if err != nil {
		return &Loader{rulesDir: "configs/rules"}
	}
	rulesDir := filepath.Join(filepath.Dir(exe), "configs", "rules")
	if _, err := os.Stat(rulesDir); err != nil {
		rulesDir = "configs/rules"
	}
	return &Loader{rulesDir: rulesDir}
}

func (l *Loader) LoadFromFile(path string) ([]*models.CleanRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules []*models.CleanRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func (l *Loader) LoadByLevel(level int) ([]*models.CleanRule, error) {
	var allRules []*models.CleanRule

	for i := 1; i <= level; i++ {
		filename := filepath.Join(l.rulesDir, fmt.Sprintf("level%d_*.json", i))
		matches, _ := filepath.Glob(filename)

		for _, match := range matches {
			rules, err := l.LoadFromFile(match)
			if err != nil {
				continue
			}
			allRules = append(allRules, rules...)
		}
	}

	return allRules, nil
}
