// Package service is
package service

import (
	"context"

	"github.com/google/uuid"
	_ "github.com/oleshko-g/oggophermart/internal/gen/accrual"
	genBalance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	_ "github.com/oleshko-g/oggophermart/internal/gen/http/accrual/client"
	genSvc "github.com/oleshko-g/oggophermart/internal/gen/service"
	genUser "github.com/oleshko-g/oggophermart/internal/gen/user"
)

type Service struct {
	Balance genBalance.Service
	User    genUser.Service
}

type JWTToken = genSvc.JWTToken

type Auther interface {
	genBalance.Auther
	UserIDFromContext(context.Context) (uuid.UUID, error)
}
