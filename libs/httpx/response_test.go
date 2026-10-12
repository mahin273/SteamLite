package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"streamlite.local/libs/errx"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := HealthHandler()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status 'ok', got %q", body["status"])
	}
}

func TestRespondError_AppError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-Id", "test-req-123")
	rr := httptest.NewRecorder()

	appErr := errx.BadRequest("INVALID_INPUT", "Username must be alphanumeric", nil)
	RespondError(rr, req, appErr)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var problem ProblemDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to parse problem JSON: %v", err)
	}

	if problem.Status != 400 {
		t.Fatalf("expected problem.Status 400, got %d", problem.Status)
	}
	if problem.Code != "INVALID_INPUT" {
		t.Fatalf("expected code INVALID_INPUT, got %q", problem.Code)
	}
	if problem.Detail != "Username must be alphanumeric" {
		t.Fatalf("expected detail 'Username must be alphanumeric', got %q", problem.Detail)
	}
	if problem.RequestID != "test-req-123" {
		t.Fatalf("expected requestId test-req-123, got %q", problem.RequestID)
	}
}

func TestRespondError_SanitizesGenericError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sensitive", nil)
	rr := httptest.NewRecorder()

	// Simulating raw internal database leak error
	sensitiveErr := errors.New("pq: connection refused to 10.0.0.5:5432 user=admin")
	RespondError(rr, req, sensitiveErr)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}

	var problem ProblemDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to parse problem JSON: %v", err)
	}

	if problem.Status != 500 {
		t.Fatalf("expected status 500, got %d", problem.Status)
	}
	if problem.Code != "INTERNAL_ERROR" {
		t.Fatalf("expected code INTERNAL_ERROR, got %q", problem.Code)
	}
	// Verify sensitive details are NOT exposed to client
	if problem.Detail != "An unexpected error occurred. Please try again later." {
		t.Fatalf("sensitive error leaked in detail: %q", problem.Detail)
	}
}
