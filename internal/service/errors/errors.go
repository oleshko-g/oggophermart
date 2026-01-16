// Package service is the package shared by all oggophermart services
package errors

import (
	"errors"

	genSvc "github.com/oleshko-g/oggophermart/internal/gen/service"
)

type svcError = genSvc.GophermartError

// New returns a pointer to a new instance of gophermart service error
func New(name string) error {
	return &svcError{Name: name}
}

var (
	// ErrInvalidInputParameter is the error value which is used to map to the 400 Bad Request HTTP Status code
	ErrInvalidInputParameter = New("Invalid input parameter")
	// ErrUserIsNotAuthenticated is the error value which is used to map to the 401 Unauthorized HTTP Status code
	ErrUserIsNotAuthenticated = New("User is not authenticated")
	// ErrInternalServiceError is the error value which is used to map to the 500 Internal Server Error HTTP Status code
	ErrInternalServiceError = New("Internal service error")
	// ErrNotImplemented is the error value which is used to map to the 501 Not Implemented HTTP Status code
	ErrNotImplemented = New("Not Implemented")

	// ErrAccrualNotStarted is the error returned when new accrual order arrives while accrual proccessing hasn't started yet
	ErrAccrualProcessingNotStarted = errors.New("Accrual processing hasn't started yet")
)
