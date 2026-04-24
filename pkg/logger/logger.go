package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Init(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	logFile, err := os.OpenFile(filepath.Join(dir, "cleaner.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w := io.MultiWriter(os.Stderr, logFile)
	log.SetOutput(w)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}
