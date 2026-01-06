// Package storage declares the type to imported by the service architecture layer
package storage

import "context"

type Storage struct {
	User    // interface
	Balance // interface
}

// User declares the storage interface for the user service
type User interface {
	RetrieveUser(id string) error
	StoreUser(ctx context.Context, login, hashedPassword string) error
}

// Balance declares the storage interfce for the balance service
type Balance interface {
	RetrieveUserBalance(userID string) (currentBalance, withdrawn int, err error)
	SaveUserTransaction(userID string, amount int) error
}
