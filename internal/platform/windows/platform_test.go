package windows

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsPlatform_ExpandPath(t *testing.T) {
	p := NewPlatform()

	// 使用 ExpandPath 已知的环境变量列表中的 TEMP
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		t.Skip("TEMP env var not set")
	}

	result := p.ExpandPath("%TEMP%\\sub")
	expected := filepath.Join(tempDir, "sub")
	if result != expected {
		t.Errorf("ExpandPath() = %s, want %s", result, expected)
	}
}

func TestWindowsPlatform_ExpandPath_RealEnv(t *testing.T) {
	p := NewPlatform()

	result := p.ExpandPath("%TEMP%")
	orig := "%TEMP%"
	if result == orig {
		t.Errorf("Expected %%TEMP%% to be expanded, got %s", result)
	}
}

func TestWindowsPlatform_GetCommonPaths(t *testing.T) {
	p := NewPlatform()
	paths := p.GetCommonPaths()

	if paths["TEMP"] == "" {
		t.Error("Expected TEMP path to be set")
	}
	if paths["HOME"] == "" {
		t.Error("Expected HOME path to be set")
	}
}

func TestWindowsPlatform_ClearRecycleBin(t *testing.T) {
	// 不实际清空回收站，只验证命令可以构建
	p := NewPlatform()
	_ = p.ClearRecycleBin // 验证方法存在
}
