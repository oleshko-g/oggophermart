// Package storage declares the type to imported by the service architecture layer
package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	User    // interface
	Balance // interface
}

// User declares the storage interface for the user service
type User interface {
	RetrieveUser(ctx context.Context, login string) (userID uuid.UUID, err error)
	StoreUser(ctx context.Context, login, hashedPassword string) error
}

// Balance declares the storage interfce for the balance service
type Balance interface {
	RetrieveUserBalance(ctx context.Context, userID uuid.UUID) (currentBalance, withdrawn int, err error)
	SaveUserTransaction(ctx context.Context, userID uuid.UUID, amount int) error
	StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, status string, createdAt time.Time) error
}
