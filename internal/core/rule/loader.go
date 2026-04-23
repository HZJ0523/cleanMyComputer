package rule

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Loader struct {
	rulesDir string
}

func NewLoader() *Loader {
	return &Loader{
		rulesDir: "configs/rules",
	}
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
		filename := filepath.Join(l.rulesDir, "level"+string(rune(i+'0'))+"_*.json")
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
