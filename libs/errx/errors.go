package errx

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
	Detail     string `json:"detail"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Detail, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Detail)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(statusCode int, code, detail string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Detail:     detail,
		Err:        err,
	}
}

func BadRequest(code, detail string, err error) *AppError {
	return New(http.StatusBadRequest, code, detail, err)
}

func Unauthorized(code, detail string, err error) *AppError {
	return New(http.StatusUnauthorized, code, detail, err)
}

func Forbidden(code, detail string, err error) *AppError {
	return New(http.StatusForbidden, code, detail, err)
}

func NotFound(code, detail string, err error) *AppError {
	return New(http.StatusNotFound, code, detail, err)
}

func Conflict(code, detail string, err error) *AppError {
	return New(http.StatusConflict, code, detail, err)
}

func Internal(detail string, err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", detail, err)
}

func AsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
