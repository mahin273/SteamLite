package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"streamlite.local/libs/errx"
)

type ProblemDetails struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Code      string `json:"code"`
	RequestID string `json:"requestId,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("failed to encode JSON response", "error", err)
		}
	}
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	reqID := r.Header.Get("X-Request-Id")

	appErr, ok := errx.AsAppError(err)
	if !ok {
		slog.Error("unhandled internal error", "error", err, "path", r.URL.Path, "requestId", reqID)
		problem := ProblemDetails{
			Type:      "about:blank",
			Title:     http.StatusText(http.StatusInternalServerError),
			Status:    http.StatusInternalServerError,
			Detail:    "An unexpected error occurred. Please try again later.",
			Code:      "INTERNAL_ERROR",
			RequestID: reqID,
		}
		RespondJSON(w, http.StatusInternalServerError, problem)
		return
	}

	if appErr.StatusCode >= http.StatusInternalServerError {
		slog.Error("internal server app error", "error", appErr.Err, "code", appErr.Code, "path", r.URL.Path, "requestId", reqID)
	}

	problem := ProblemDetails{
		Type:      "about:blank",
		Title:     http.StatusText(appErr.StatusCode),
		Status:    appErr.StatusCode,
		Detail:    appErr.Detail,
		Code:      appErr.Code,
		RequestID: reqID,
	}

	RespondJSON(w, appErr.StatusCode, problem)
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}

func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondError(w, r, errx.NotFound("ROUTE_NOT_FOUND", "The requested endpoint does not exist", nil))
	}
}
