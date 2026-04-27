package cleaner

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	absBase, _ := filepath.Abs(baseDir)
	return &QuarantineManager{
		baseDir:        absBase,
		retentionHours: 24,
	}, nil
}

func (q *QuarantineManager) Quarantine(srcPath string) error {
	absSrc, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	if strings.HasPrefix(absSrc, q.baseDir+string(os.PathSeparator)) {
		return fmt.Errorf("cannot quarantine a file already in quarantine directory")
	}

	fileName := filepath.Base(srcPath)
	quarantineName := fmt.Sprintf("%d_%d_%s", time.Now().UnixNano(), rand.Intn(10000), fileName)
	dstPath := filepath.Join(q.baseDir, quarantineName)

	info, err := os.Stat(srcPath)
	var size int64
	if err == nil {
		size = info.Size()
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return err
	}

	if q.OnQuarantine != nil {
		now := time.Now()
		if err := q.OnQuarantine(QuarantineRecord{
			OriginalPath:   srcPath,
			QuarantinePath: dstPath,
			Size:           size,
			RiskScore:      0,
			CreatedAt:      now,
			ExpiresAt:      now.Add(time.Duration(q.retentionHours) * time.Hour),
		}); err != nil {
			return fmt.Errorf("quarantine succeeded but failed to persist record: %w", err)
		}
	}

	return nil
}

func (q *QuarantineManager) Restore(quarantinePath, originalPath string) error {
	if err := os.MkdirAll(filepath.Dir(originalPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	return os.Rename(quarantinePath, originalPath)
}
