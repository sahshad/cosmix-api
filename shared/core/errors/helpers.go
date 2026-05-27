package errors

import "net/http"

func NewBadRequest(field string, message string) *AppError {
	return &AppError{
		Code:    CodeBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
		Field:   field,
	}
}

func NewUnauthorized(message string) *AppError {
	if message == "" {
		message = "unauthorized"
	}

	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbidden(message string) *AppError {
	if message == "" {
		message = "forbidden"
	}

	return &AppError{
		Code:    CodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: resource + " not found",
		Status:  http.StatusNotFound,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

func NewValidation(field string, message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
		Status:  http.StatusBadRequest,
		Field:   field,
	}
}

func NewInternal(err error) *AppError {
	return &AppError{
		Code:    CodeInternalServerError,
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}
