package analyzer

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type DuplicateGroup struct {
	Hash  string
	Size  int64
	Paths []string
}

type DuplicateFinder struct {
	minSize int64
}

func NewDuplicateFinder(minSize int64) *DuplicateFinder {
	if minSize <= 0 {
		minSize = 1024 // 1KB minimum
	}
	return &DuplicateFinder{minSize: minSize}
}

func (d *DuplicateFinder) FindDuplicates(root string) ([]DuplicateGroup, error) {
	sizeMap := make(map[int64][]string)

	filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		fi, err := info.Info()
		if err != nil || fi.Size() < d.minSize {
			return nil
		}
		sizeMap[fi.Size()] = append(sizeMap[fi.Size()], path)
		return nil
	})

	var candidates [][]string
	for _, paths := range sizeMap {
		if len(paths) > 1 {
			candidates = append(candidates, paths)
		}
	}

	hashMap := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, group := range candidates {
		for _, path := range group {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				h, err := fileHash(p)
				if err != nil {
					return
				}
				mu.Lock()
				hashMap[h] = append(hashMap[h], p)
				mu.Unlock()
			}(path)
		}
	}
	wg.Wait()

	var result []DuplicateGroup
	for hash, paths := range hashMap {
		if len(paths) > 1 {
			var size int64
			if fi, err := os.Stat(paths[0]); err == nil {
				size = fi.Size()
			}
			result = append(result, DuplicateGroup{Hash: hash, Size: size, Paths: paths})
		}
	}
	return result, nil
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}
