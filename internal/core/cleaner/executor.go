package cleaner

import (
	"context"
	"os"
	"os/exec"
	"strings"
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
	Cleaned      []string
	Quarantined  []string
	Failed       []string
	FreedSize    int64
	QuarantinedSize int64
	StartTime    time.Time
	EndTime      time.Time
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
		action, err := e.cleanFile(file)
		if err != nil {
			result.Failed = append(result.Failed, file.Path)
			continue
		}
		switch action {
		case "deleted":
			result.Cleaned = append(result.Cleaned, file.Path)
			result.FreedSize += file.Size
		case "quarantined":
			result.Quarantined = append(result.Quarantined, file.Path)
			result.QuarantinedSize += file.Size
			result.Cleaned = append(result.Cleaned, file.Path)
		case "command":
			result.Cleaned = append(result.Cleaned, file.Path)
		}
	}
	result.EndTime = time.Now()
	return result, nil
}

func (e *Executor) cleanFile(file *FileItem) (string, error) {
	if e.dryRun {
		return "deleted", nil
	}
	if !strings.Contains(file.Path, "\\") && !strings.Contains(file.Path, "/") {
		return "command", e.executeCommand(file.Path)
	}
	if file.RiskScore > 60 {
		return "quarantined", e.quarantine.Quarantine(file.Path)
	}
	return "deleted", e.deleteWithRetry(file.Path)
}

const maxRetries = 2

func (e *Executor) deleteWithRetry(path string) error {
	for i := 0; i < maxRetries; i++ {
		err := os.Remove(path)
		if err == nil {
			return nil
		}
		if os.IsPermission(err) {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return os.Remove(path)
}

func (e *Executor) executeCommand(cmd string) error {
	if cmd == "Clear-RecycleBin" {
		return exec.Command("powershell", "-Command", "Clear-RecycleBin -Force").Run()
	}
	return exec.Command("cmd", "/C", cmd).Run()
}
