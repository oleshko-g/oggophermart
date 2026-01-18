package accrual

type AccrualError struct {
	// identifier to map an error to HTTP status codes
	RetryAfter int
	Message    string
}

// Error returns an error Message.
func (e *AccrualError) Error() string {
	return e.Message
}
