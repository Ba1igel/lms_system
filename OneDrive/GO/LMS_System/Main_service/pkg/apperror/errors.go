package apperror

import (
	"fmt"
	"net/http"
)

// AppError — типизированная ошибка с HTTP-кодом.
// Создаётся в сервисном слое, рендерится в ErrorHandlerMiddleware.
type AppError struct {
	Code    int    // HTTP status code
	Message string // публичное сообщение для клиента
	Err     error  // оригинальная причина (только для логов, не отдаётся клиенту)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap позволяет errors.As/errors.Is проходить через цепочку.
func (e *AppError) Unwrap() error {
	return e.Err
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func BadRequest(msg string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}
