// Package service is the package shared by all oggophermart services
package errors

import (
	genSvc "github.com/oleshko-g/oggophermart/internal/gen/service"
	goa "goa.design/goa/v3/pkg"
)

type svcError struct {
	genSvc.GophermartError
}

var _ goa.GoaErrorNamer= (*svcError)(nil)

// Error return the Name field of the wrapped [GophermartError], because the standard generated code return an empty string
func (ogerr svcError) Error() string {
	return ogerr.Name
}

// String return the Name field of the wrapped [GophermartError]
func (ogerr svcError) String() string {
	return ogerr.Name
}

// New returns a pointer to a new instance of gophermart service error
func New(name string) error {
	return &svcError{genSvc.GophermartError{Name: name}}
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
)
