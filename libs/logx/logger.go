package logx

import (
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

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}).WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
		slog.String("env", env),
	})

	return slog.New(handler)
}
