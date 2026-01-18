// Package storage declares the type to imported by the service architecture layer
package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	genDBSQL "github.com/oleshko-g/oggophermart/internal/gen/storage/db/sql"
)

type Storager interface {
	User    // interface
	Balance // interface
}

type Storage struct {
	Transacter //interface
	User       // interface
	Balance    // interface
}

type Order = genDBSQL.Order

// User declares the storage interface for the user service
type User interface {
	RetrieveUser(ctx context.Context, login string) (userID uuid.UUID, err error)
	RetreiveUserPassword(ctx context.Context, login string) (hashedPassword string, err error)
	StoreUser(ctx context.Context, login, hashedPassword string) error
}

// Balance declares the storage interfce for the balance service
type Balance interface {
	Transacter
	RetrieveUserBalance(ctx context.Context, userID uuid.UUID) (currentBalance, withdrawn int, err error)
	SaveUserTransaction(ctx context.Context, userID uuid.UUID, amount int) error
	StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, status string, createdAt time.Time) error
	RetreiveOrderUser(ctx context.Context, orderNumber string) (userID uuid.UUID, err error)
	RetrieaveUserOrders(ctx context.Context, userID uuid.UUID) ([]genDBSQL.SelectOrdersByUserIDRow, error)
	Retrieve(ctx context.Context, userID uuid.UUID) (genDBSQL.SelectBalanceByUserIDRow, error)
	RetrieveOrderIDsForAccrual(ctx context.Context) ([]uuid.UUID, error)
	RetrieveOrderForAccrual(ctx context.Context, orderID uuid.UUID) (Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
	StoreUserAccrual(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, amount int32) error
}

type Transaction interface {
	Commit() error
	Rollback() error
}

type Tx struct {
	Tx Transaction
	Balance
}

type Transacter interface {
	BeginTx(context.Context) (*Tx, error)
}
