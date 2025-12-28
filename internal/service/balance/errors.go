package oggophermart

import(
		"github.com/oleshko-g/oggophermart/internal/service/errors"
)

var (
	OwnerMismatch = errors.New("The order belongs to another user")
)
