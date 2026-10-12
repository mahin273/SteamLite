package logx

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func New(serviceName, env string) *slog.Logger {
	var level slog.Level
	if strings.ToLower(env) == "production" || strings.ToLower(env) == "prod" {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}

	var output io.Writer = os.Stdout
	if fileWriter, err := newDailyFileWriter(logDir, serviceName); err == nil {
		output = io.MultiWriter(os.Stdout, fileWriter)
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	}).WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
		slog.String("env", env),
	})

	return slog.New(handler)
}
