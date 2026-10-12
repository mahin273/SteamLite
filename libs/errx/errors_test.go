package errx

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Formatting(t *testing.T) {
	underlying := errors.New("db connection lost")
	appErr := New(http.StatusInternalServerError, "DB_ERROR", "Failed to query database", underlying)

	if appErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, appErr.StatusCode)
	}

	if appErr.Code != "DB_ERROR" {
		t.Fatalf("expected code DB_ERROR, got %s", appErr.Code)
	}

	if !errors.Is(appErr, underlying) {
		t.Fatalf("expected error unwrapping to match underlying error")
	}

	expectedMsg := "[DB_ERROR] Failed to query database: db connection lost"
	if appErr.Error() != expectedMsg {
		t.Fatalf("expected error message %q, got %q", expectedMsg, appErr.Error())
	}
}

func TestAsAppError(t *testing.T) {
	appErr := BadRequest("INVALID_EMAIL", "Email is not valid", nil)
	extracted, ok := AsAppError(appErr)
	if !ok {
		t.Fatalf("expected ok to be true for AppError")
	}
	if extracted.Code != "INVALID_EMAIL" {
		t.Fatalf("expected code INVALID_EMAIL, got %s", extracted.Code)
	}

	genericErr := errors.New("plain error")
	_, okGeneric := AsAppError(genericErr)
	if okGeneric {
		t.Fatalf("expected ok to be false for plain generic error")
	}
}
