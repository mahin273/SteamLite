package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDailyFileWriter_CreatesAndWritesLog(t *testing.T) {
	tempDir := t.TempDir()
	serviceName := "test-svc"

	w, err := newDailyFileWriter(tempDir, serviceName)
	if err != nil {
		t.Fatalf("failed to create daily file writer: %v", err)
	}
	defer func() {
		if closer, ok := w.(interface{ Close() error }); ok {
			_ = closer.Close()
		}
	}()

	msg := []byte("{\"level\":\"INFO\",\"msg\":\"test log message\"}\n")
	n, err := w.Write(msg)
	if err != nil {
		t.Fatalf("failed to write log: %v", err)
	}
	if n != len(msg) {
		t.Fatalf("expected write %d bytes, got %d", len(msg), n)
	}

	today := time.Now().Format("02_01_2006")
	expectedPath := filepath.Join(tempDir, serviceName, fmt.Sprintf("%s.log", today))

	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read created log file %s: %v", expectedPath, err)
	}

	if !strings.Contains(string(content), "test log message") {
		t.Fatalf("expected log file to contain 'test log message', got %s", string(content))
	}
}
