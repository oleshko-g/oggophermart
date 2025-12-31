// Package balance declares the user storage interface
package storage

type Balance interface {
	RetrieveUserBalance(userID string) (currentBalance, withdrawn int, err error)
	SaveUserTransaction(userID string, amount int) (error)
}
