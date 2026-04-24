package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hzj0523/cleanMyComputer/internal/models"
	"github.com/hzj0523/cleanMyComputer/pkg/i18n"
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
	sb.WriteString(i18n.T("report.title") + "\n\n")
	sb.WriteString(fmt.Sprintf(i18n.T("report.scan_level")+"\n", report.ScanLevel))
	sb.WriteString(fmt.Sprintf(i18n.T("report.scan_time")+"\n",
		report.StartTime.Format("2006-01-02 15:04:05"),
		report.EndTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf(i18n.T("report.item_count")+"\n", len(report.Items)))
	sb.WriteString(fmt.Sprintf(i18n.T("report.freed_space")+"\n", formatBytes(report.FreedSize)))
	sb.WriteString(fmt.Sprintf(i18n.T("report.failed_count")+"\n\n", report.Failed))

	sb.WriteString(i18n.T("report.detail_header") + "\n")
	for i, item := range report.Items {
		risk := i18n.T("risk.safe")
		if item.RiskScore > 60 {
			risk = i18n.T("risk.high")
		} else if item.RiskScore > 30 {
			risk = i18n.T("risk.moderate")
		}
		sb.WriteString(fmt.Sprintf(i18n.T("report.detail_item")+"\n",
			i+1, sanitizePath(item.Path),
			i18n.T("label.size"), formatBytes(item.Size),
			i18n.T("label.risk"), risk))
	}
	sb.WriteString("\n" + i18n.T("report.footer") + "\n")
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
