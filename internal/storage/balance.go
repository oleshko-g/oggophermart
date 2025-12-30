// Package balance declares the user storage interface
package storage

type Balance interface {
	RetrieveUserBalance(userID string) (int, error)
	SaveUserTransaction(userID string, amount int) (error)
}
