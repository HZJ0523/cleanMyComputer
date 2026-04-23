package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type QuarantineManager struct {
	baseDir string
}

func NewQuarantineManager(baseDir string) *QuarantineManager {
	os.MkdirAll(baseDir, 0755)
	return &QuarantineManager{baseDir: baseDir}
}

func (q *QuarantineManager) Quarantine(srcPath string) error {
	fileName := filepath.Base(srcPath)
	quarantineName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
	dstPath := filepath.Join(q.baseDir, quarantineName)
	return os.Rename(srcPath, dstPath)
}

func (q *QuarantineManager) Restore(quarantinePath, originalPath string) error {
	os.MkdirAll(filepath.Dir(originalPath), 0755)
	return os.Rename(quarantinePath, originalPath)
}
