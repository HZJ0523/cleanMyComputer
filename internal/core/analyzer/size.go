package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type LargeFile struct {
	Path string
	Size int64
}

type LargeFileFinder struct {
	threshold  int64
	maxResults int
}

func NewLargeFileFinder(threshold int64) *LargeFileFinder {
	if threshold <= 0 {
		threshold = 100 * 1024 * 1024
	}
	return &LargeFileFinder{threshold: threshold, maxResults: 50}
}

func (l *LargeFileFinder) FindLargeFiles(root string) ([]LargeFile, error) {
	var files []LargeFile

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Size() >= l.threshold {
			files = append(files, LargeFile{Path: path, Size: info.Size()})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk %s: %w", root, err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > l.maxResults {
		files = files[:l.maxResults]
	}
	return files, nil
}
