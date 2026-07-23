package domain

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrKeyNotFound        ErrorCode = "KEY_NOT_FOUND"
	ErrKeyInactive        ErrorCode = "KEY_INACTIVE"
	ErrTenantNotFound     ErrorCode = "TENANT_NOT_FOUND"
	ErrInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrHSMOperationFailed ErrorCode = "HSM_OPERATION_FAILED"
	ErrK8sSecretNotFound  ErrorCode = "K8S_SECRET_NOT_FOUND"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
	Details map[string]interface{} // interno: secretName, namespace, paths, etc.
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error { return e.Err }

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

func NewDomainError(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}
