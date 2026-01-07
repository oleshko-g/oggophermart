package user

import (
	"github.com/oleshko-g/oggophermart/internal/service/errors"
)

var (
	// ErrLoginTaken is returnd when a registration failed dew to a taken login
	ErrLoginTaken = errors.New("Login is taken already")
)
