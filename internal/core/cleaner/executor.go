package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var allowedCommands = map[string][]string{
	"Clear-RecycleBin": {"powershell", "-Command", "Clear-RecycleBin -Force"},
}

type Executor struct {
	dryRun bool
}

type CleanTask struct {
	Files     []*FileItem
	TotalSize int64
}

type FileItem struct {
	Path      string
	Size      int64
	RiskScore int
	Type      string
}

type CleanResult struct {
	Cleaned   []string
	Failed    []string
	FreedSize int64
	StartTime time.Time
	EndTime   time.Time
}

func NewExecutor() *Executor {
	return &Executor{dryRun: false}
}

func (e *Executor) Execute(ctx context.Context, task *CleanTask) (*CleanResult, error) {
	result := &CleanResult{StartTime: time.Now()}
	for _, file := range task.Files {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		if err := e.cleanFile(ctx, file); err != nil {
			result.Failed = append(result.Failed, file.Path)
			continue
		}
		result.Cleaned = append(result.Cleaned, file.Path)
		result.FreedSize += file.Size
	}
	result.EndTime = time.Now()
	return result, nil
}

func (e *Executor) cleanFile(ctx context.Context, file *FileItem) error {
	if e.dryRun {
		return nil
	}
	if file.Type == "command" {
		return e.executeCommand(ctx, file.Path)
	}
	return e.deleteWithRetry(file.Path)
}

const maxRetries = 2

func (e *Executor) deleteWithRetry(path string) error {
	for i := 0; i < maxRetries; i++ {
		err := os.Remove(path)
		if err == nil {
			return nil
		}
		if os.IsNotExist(err) || os.IsPermission(err) {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return os.Remove(path)
}

func (e *Executor) executeCommand(ctx context.Context, cmd string) error {
	args, ok := allowedCommands[cmd]
	if !ok {
		return fmt.Errorf("unrecognized command: %q", cmd)
	}
	return exec.CommandContext(ctx, args[0], args[1:]...).Run()
}

func IsCommandTarget(path string) bool {
	return !strings.Contains(path, "\\") && !strings.Contains(path, "/")
}
