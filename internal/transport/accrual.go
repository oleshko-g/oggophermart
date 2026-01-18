package transport

import "context"

// Accrual is the interface to accrual system
type Accrual interface {
	// FetchOrderAccrual method
	FetchOrderAccrual(ctx context.Context, payload FetchOrderAccrualPayload) (*FetchOrderAccrualResult, error)
}

// FetchOrderAccrualPayload is the payload type of the accrual system of [FetchOrderAccrual] method
// FetchOrderAccrual method
type FetchOrderAccrualPayload struct {
	Number string
}

// FetchOrderAccrualResult is the result type of the accrual system of [FetchOrderAccrual] method
type FetchOrderAccrualResult struct {
	Order   string
	Status  OrderAccrualStatus
	Accrual *float64
}

type AccrualError struct {
	// identifier to map an error to HTTP status codes
	RetryAfter int
	Message    string
}

// Error returns an error Message.
func (e *AccrualError) Error() string {
	return e.Message
}

type OrderAccrualStatus string

const (
	OrderAccrualStatusRegistered = "REGISTERED"
	OrderAccrualStatusProcessing = "PROCESSING"
	OrderAccrualStatusProcessed  = "PROCESSED"
	OrderAccrualStatusInvalid    = "INVALID"
)
