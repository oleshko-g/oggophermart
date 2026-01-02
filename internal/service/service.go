package service

import (
	_ "github.com/oleshko-g/oggophermart/internal/gen/accrual"
	"github.com/oleshko-g/oggophermart/internal/gen/balance"
	_ "github.com/oleshko-g/oggophermart/internal/gen/http/accrual/client"
	"github.com/oleshko-g/oggophermart/internal/gen/user"
)

type Service struct {
	Balance balance.Service
	User    user.Service
}
