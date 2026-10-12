package http

import (
	"log/slog"
	"net/http"

	"streamlite.local/libs/httpx"
)

func NewRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", httpx.HealthHandler())

	var handler http.Handler = mux
	handler = httpx.RequestLogger(logger)(handler)
	handler = httpx.Recovery(logger)(handler)
	handler = httpx.RequestID(handler)

	return handler
}
