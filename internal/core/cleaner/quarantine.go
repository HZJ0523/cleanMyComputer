package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// QuarantineRecord 记录隔离信息，由调用者持久化
type QuarantineRecord struct {
	OriginalPath   string
	QuarantinePath string
	Size           int64
	RiskScore      int
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

type QuarantineManager struct {
	baseDir        string
	retentionHours int
	OnQuarantine   func(record QuarantineRecord) error
}

func NewQuarantineManager(baseDir string) (*QuarantineManager, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create quarantine directory: %w", err)
	}
	return &QuarantineManager{
		baseDir:        baseDir,
		retentionHours: 24,
	}, nil
}

func (q *QuarantineManager) SetRetentionHours(hours int) {
	q.retentionHours = hours
}

func (q *QuarantineManager) Quarantine(srcPath string) error {
	fileName := filepath.Base(srcPath)
	quarantineName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
	dstPath := filepath.Join(q.baseDir, quarantineName)

	// Get file info before moving
	info, err := os.Stat(srcPath)
	var size int64
	if err == nil {
		size = info.Size()
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return err
	}

	// Notify caller for persistence
	if q.OnQuarantine != nil {
		now := time.Now()
		q.OnQuarantine(QuarantineRecord{
			OriginalPath:   srcPath,
			QuarantinePath: dstPath,
			Size:           size,
			CreatedAt:      now,
			ExpiresAt:      now.Add(time.Duration(q.retentionHours) * time.Hour),
		})
	}

	return nil
}

func (q *QuarantineManager) Restore(quarantinePath, originalPath string) error {
	os.MkdirAll(filepath.Dir(originalPath), 0755)
	return os.Rename(quarantinePath, originalPath)
}
