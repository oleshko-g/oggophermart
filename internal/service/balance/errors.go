package balance

import (
	"github.com/oleshko-g/oggophermart/internal/service/errors"
)

var (
	// ErrOwnerMismatch is the error value which is used to map to the 409 Conflict HTTP Status code
	ErrOwnerMismatch = errors.New("The order belongs to another user")
	// ErrInvalidOrderNumber is the error value which is used to map to the 422 Unprocessable Entity HTTP Status code
	ErrInvalidOrderNumber = errors.New("Invalid order number")
)
