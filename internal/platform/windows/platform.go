package windows

import (
	"errors"
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsPlatform struct{}

func NewPlatform() *WindowsPlatform {
	return &WindowsPlatform{}
}

type DiskUsage struct {
	TotalGB float64
	UsedGB  float64
	FreeGB  float64
}

func (w *WindowsPlatform) GetDiskUsage(path string) (*DiskUsage, error) {
	var freeAvailableBytes uint64
	var totalBytes uint64
	var freeBytes uint64

	root := filepath.VolumeName(path) + "\\"
	kernel32 := windows.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")
	rootPtr, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return nil, fmt.Errorf("invalid root path %q: %w", root, err)
	}

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(rootPtr)),
		uintptr(unsafe.Pointer(&freeAvailableBytes)),
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
