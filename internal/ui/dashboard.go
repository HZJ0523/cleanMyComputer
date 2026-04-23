package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/sys/windows"
)

func (a *App) newDashboard() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("电脑垃圾清理工具", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Show disk space
	var diskInfo string
	homedir, _ := os.UserHomeDir()
	if usage, err := getDiskUsage(homedir); err == nil {
		diskInfo = fmt.Sprintf("系统盘: 已用 %.1f GB / 总计 %.1f GB (可用 %.1f GB)",
			usage.UsedGB, usage.TotalGB, usage.FreeGB)
	}

	diskLabel := widget.NewLabel(diskInfo)
	statusLabel := widget.NewLabel("选择扫描模式开始清理")
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	quickScan := widget.NewButton("快速扫描", func() {
		statusLabel.SetText("正在快速扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(1)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			count := len(a.state.ScanItems)
			statusLabel.SetText(fmt.Sprintf("扫描完成，发现 %d 个可清理项", count))
			progressBar.Hide()
			dialog.ShowInformation("扫描完成",
				fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", count), a.window)
			a.selectTab(1)
		}()
	})

	fullScan := widget.NewButton("完整扫描", func() {
		statusLabel.SetText("正在完整扫描...")
		progressBar.Show()
		go func() {
			err := a.state.RunScan(3)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			count := len(a.state.ScanItems)
			statusLabel.SetText(fmt.Sprintf("扫描完成，发现 %d 个可清理项", count))
			progressBar.Hide()
			dialog.ShowInformation("扫描完成",
				fmt.Sprintf("发现 %d 个可清理项，请查看扫描结果", count), a.window)
			a.selectTab(1)
		}()
	})

	return container.NewVBox(
		title,
		diskLabel,
		statusLabel,
		progressBar,
		container.NewHBox(quickScan, fullScan),
	)
}

type diskUsage struct {
	TotalGB float64
	UsedGB  float64
	FreeGB  float64
}

func getDiskUsage(path string) (*diskUsage, error) {
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
		return nil, fmt.Errorf("GetDiskFreeSpaceEx failed")
	}

	usedGB := float64(totalBytes-freeBytes) / 1024 / 1024 / 1024
	return &diskUsage{
		TotalGB: float64(totalBytes) / 1024 / 1024 / 1024,
		UsedGB:  usedGB,
		FreeGB:  float64(freeBytes) / 1024 / 1024 / 1024,
	}, nil
}
