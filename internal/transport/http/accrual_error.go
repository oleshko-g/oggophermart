package http

type AccrualError struct {
	// identifier to map an error to HTTP status codes
	Name       string `json:"-"`
	RetryAfter int
	Message    string
}

// Error returns an error description.
func (e *AccrualError) Error() string {
	return ""
}

// ErrorName returns "AccrualError".
//
// Deprecated: Use GoaErrorName - https://github.com/goadesign/goa/issues/3105
func (e *AccrualError) ErrorName() string {
	return e.GoaErrorName()
}

// GoaErrorName returns "AccrualError".
func (e *AccrualError) GoaErrorName() string {
	return e.Name
}
