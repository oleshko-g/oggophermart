package balance

import (
	"github.com/oleshko-g/oggophermart/internal/service/errors"
)

var (
	// OwnerMismatch is the error value which is used to map to the 409 Conflict HTTP Status code
	OwnerMismatch = errors.New("The order belongs to another user")
	// InvalidOrderNumber is the error value which is used to map to the 422 Unprocessable Entity HTTP Status code
	InvalidOrderNumber = errors.New("Invalid order number")
)
