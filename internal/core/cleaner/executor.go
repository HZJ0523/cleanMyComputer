package cleaner

import (
	"context"
	"os"
	"time"
)

type Executor struct {
	quarantine *QuarantineManager
	dryRun     bool
}

type CleanTask struct {
	Files     []*FileItem
	TotalSize int64
}

type FileItem struct {
	Path      string
	Size      int64
	RiskScore int
}

type CleanResult struct {
	Cleaned   []string
	Failed    []string
	FreedSize int64
	StartTime time.Time
	EndTime   time.Time
}

func NewExecutor(qm *QuarantineManager) *Executor {
	return &Executor{quarantine: qm, dryRun: false}
}

func (e *Executor) Execute(ctx context.Context, task *CleanTask) (*CleanResult, error) {
	result := &CleanResult{StartTime: time.Now()}
	for _, file := range task.Files {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		if err := e.cleanFile(file); err != nil {
			result.Failed = append(result.Failed, file.Path)
			continue
		}
		result.Cleaned = append(result.Cleaned, file.Path)
		result.FreedSize += file.Size
	}
	result.EndTime = time.Now()
	return result, nil
}

func (e *Executor) cleanFile(file *FileItem) error {
	if e.dryRun {
		return nil
	}
	if file.RiskScore > 60 {
		return e.quarantine.Quarantine(file.Path)
	}
	return os.Remove(file.Path)
}
