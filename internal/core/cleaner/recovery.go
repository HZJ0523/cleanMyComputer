package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
)

type Recovery struct {
	qm *QuarantineManager
}

func NewRecovery(qm *QuarantineManager) *Recovery {
	return &Recovery{qm: qm}
}

func (r *Recovery) RestoreFile(quarantinePath, originalPath string) error {
	if r.qm != nil {
		return r.qm.Restore(quarantinePath, originalPath)
	}
	if err := os.MkdirAll(filepath.Dir(originalPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	return os.Rename(quarantinePath, originalPath)
}
