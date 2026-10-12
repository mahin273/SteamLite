package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type dailyFileWriter struct {
	mu         sync.Mutex
	dir        string
	currentDay string
	file       *os.File
}

func newDailyFileWriter(baseDir, serviceName string) (io.Writer, error) {
	serviceLogDir := filepath.Join(baseDir, serviceName)
	if err := os.MkdirAll(serviceLogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %q: %w", serviceLogDir, err)
	}

	w := &dailyFileWriter{
		dir: serviceLogDir,
	}

	today := time.Now().Format("02_01_2006")
	if err := w.rotate(today); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *dailyFileWriter) rotate(today string) error {
	if w.file != nil {
		_ = w.file.Close()
	}

	filename := filepath.Join(w.dir, fmt.Sprintf("%s.log", today))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %w", filename, err)
	}

	w.file = f
	w.currentDay = today
	return nil
}

func (w *dailyFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("02_01_2006")
	if today != w.currentDay || w.file == nil {
		if err := w.rotate(today); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

func (w *dailyFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}
