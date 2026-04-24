package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
)

type Report struct {
	ScanLevel  int
	Items      []*models.ScanItem
	StartTime  time.Time
	EndTime    time.Time
	FreedSize  int64
	Failed     int
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateText(report *Report) string {
	var sb strings.Builder
	sb.WriteString("=== 清理报告 ===\n\n")
	sb.WriteString(fmt.Sprintf("扫描级别: %d\n", report.ScanLevel))
	sb.WriteString(fmt.Sprintf("扫描时间: %s - %s\n",
		report.StartTime.Format("2006-01-02 15:04:05"),
		report.EndTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("扫描项目数: %d\n", len(report.Items)))
	sb.WriteString(fmt.Sprintf("释放空间: %s\n", formatBytes(report.FreedSize)))
	sb.WriteString(fmt.Sprintf("失败数: %d\n\n", report.Failed))

	sb.WriteString("--- 详细列表 ---\n")
	for i, item := range report.Items {
		risk := "安全"
		if item.RiskScore > 60 {
			risk = "高风险"
		} else if item.RiskScore > 30 {
			risk = "中等风险"
		}
		sb.WriteString(fmt.Sprintf("%d. %s (大小: %s, 风险: %s)\n",
			i+1, sanitizePath(item.Path), formatBytes(item.Size), risk))
	}
	sb.WriteString("\n--- 报告结束 ---\n")
	return sb.String()
}

func (g *Generator) ExportToFile(report *Report, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	content := g.GenerateText(report)
	return os.WriteFile(path, []byte(content), 0644)
}

func sanitizePath(path string) string {
	username := os.Getenv("USERNAME")
	if username != "" {
		path = strings.ReplaceAll(path, username, "<USER>")
	}
	return path
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
