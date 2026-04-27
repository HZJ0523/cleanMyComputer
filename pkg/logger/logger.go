package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

func Init(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	logFile, err = os.OpenFile(filepath.Join(dir, "cleaner.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w := io.MultiWriter(os.Stderr, logFile)
	log.SetOutput(w)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}

func Close() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}
