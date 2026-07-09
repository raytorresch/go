package exceptions

import (
	"errors"
	"fmt"
)

type ErrCode string

const (
	ErrCodeNotFound     ErrCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrCode = "UNAUTHORIZED"
	ErrCodeInternal     ErrCode = "INTERNAL_ERROR"
)

type DomainError struct {
	Code    ErrCode
	Message string
	Err     error
	Details map[string]interface{}
}

func (e *DomainError) Unwrap() error { return e.Err }

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) WithDetails(key string, value interface{}) *DomainError {
	e.Details[key] = value
	return e
}

func GetErrorCode(err error) string {
	if err == nil {
		return "SUCCESS"
	}

	var de *DomainError

	if errors.As(err, &de) {
		return string(de.Code)
	}

	return "INTERNAL_ERROR"
}

func NewDomainError(code ErrCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}
