package windows

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type WindowsPlatform struct{}

func NewPlatform() *WindowsPlatform {
	return &WindowsPlatform{}
}

func (w *WindowsPlatform) ExpandPath(path string) string {
	// 处理 Windows %VAR% 格式
	result := os.ExpandEnv(path)
	// os.ExpandEnv 不处理 %VAR%，手动处理
	if strings.Contains(result, "%") {
		for _, envVar := range []string{
			"TEMP", "TMP", "APPDATA", "LOCALAPPDATA",
			"USERPROFILE", "PROGRAMFILES", "PROGRAMFILES(X86)",
			"SYSTEMROOT", "WINDIR", "HOMEDRIVE", "HOMEPATH",
		} {
			value := os.Getenv(envVar)
			if value != "" {
				result = strings.ReplaceAll(result, "%"+envVar+"%", value)
				result = strings.ReplaceAll(result, "%"+strings.ToLower(envVar)+"%", value)
			}
		}
	}
	return filepath.Clean(result)
}

func (w *WindowsPlatform) IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func (w *WindowsPlatform) ClearRecycleBin() error {
	// 使用 PowerShell 清空回收站
	cmd := exec.Command("powershell", "-Command", "Clear-RecycleBin -Force")
	return cmd.Run()
}

func (w *WindowsPlatform) GetCommonPaths() map[string]string {
	home, _ := os.UserHomeDir()
	return map[string]string{
		"HOME":         home,
		"TEMP":         os.Getenv("TEMP"),
		"APPDATA":      os.Getenv("APPDATA"),
		"LOCALAPPDATA": os.Getenv("LOCALAPPDATA"),
		"PROGRAMFILES": os.Getenv("ProgramFiles"),
		"SYSTEMROOT":   os.Getenv("SystemRoot"),
	}
}
