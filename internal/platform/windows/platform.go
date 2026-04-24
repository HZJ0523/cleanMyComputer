package windows

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsPlatform struct{}

func NewPlatform() *WindowsPlatform {
	return &WindowsPlatform{}
}

func (w *WindowsPlatform) ExpandPath(path string) string {
	result := os.ExpandEnv(path)
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
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

func (w *WindowsPlatform) ClearRecycleBin() error {
	return exec.Command("powershell", "-Command", "Clear-RecycleBin -Force").Run()
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

type DiskUsage struct {
	TotalGB float64
	UsedGB  float64
	FreeGB  float64
}

func (w *WindowsPlatform) GetDiskUsage(path string) (*DiskUsage, error) {
	var freeBytes uint64
	var totalBytes uint64

	root := filepath.VolumeName(path) + "\\"
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")
	rootPtr, _ := windows.UTF16PtrFromString(root)

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(rootPtr)),
		0,
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&freeBytes)),
	)
	if ret == 0 {
		return nil, ErrGetDiskFreeSpaceFailed
	}

	return &DiskUsage{
		TotalGB: float64(totalBytes) / 1024 / 1024 / 1024,
		UsedGB:  float64(totalBytes-freeBytes) / 1024 / 1024 / 1024,
		FreeGB:  float64(freeBytes) / 1024 / 1024 / 1024,
	}, nil
}

var ErrGetDiskFreeSpaceFailed = errors.New("GetDiskFreeSpaceEx failed")
