// Package service is
package service

import (
	"context"

	"github.com/google/uuid"
	_ "github.com/oleshko-g/oggophermart/internal/gen/accrual"
	genBalance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	genAccrualHTTPClient "github.com/oleshko-g/oggophermart/internal/gen/http/accrual/client"
	genUser "github.com/oleshko-g/oggophermart/internal/gen/user"
)

type Service struct {
	User genUser.Service
	Balance
}

type Auther interface {
	genBalance.Auther
	UserIDFromContext(context.Context) (uuid.UUID, error)
}

type Balance interface {
	genBalance.Service
	ProcessAccruals(ctx context.Context, client genAccrualHTTPClient.Client) error
}
