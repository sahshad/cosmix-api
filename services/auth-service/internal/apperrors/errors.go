package apperrors

import "fmt"

type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e *NotFoundError) Error() string {
	if e.ID != nil {
		return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "unauthorized"
}

type BadRequestError struct {
	Field   string
	Message string
}

func (e *BadRequestError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "internal server error"
}

func NewNotFound(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewUnauthorized(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

func NewBadRequest(field, message string) *BadRequestError {
	return &BadRequestError{Field: field, Message: message}
}

func NewInternal(message string) *InternalServerError {
	return &InternalServerError{Message: message}
}