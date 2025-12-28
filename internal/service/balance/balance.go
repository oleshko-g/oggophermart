// Package balance is the implementation of [balance]
package oggophermart

import (
	"context"

	balance "github.com/oleshko-g/oggophermart/internal/gen/balance"
	"goa.design/clue/log"
)

// balance service example implementation.
// The example methods log the requests and return zero values.
type balanceSvc struct{}

var _ balance.Service = (*balanceSvc)(nil)

// NewBalance returns the balance service implementation.
func NewBalance() balance.Service {
	return &balanceSvc{}
}

// PostOrder implements post order.
func (s *balanceSvc) PostOrder(ctx context.Context) (res *balance.PostOrderResult, err error) {
	res = &balance.PostOrderResult{}
	log.Printf(ctx, "balance.post order")
	return
}
