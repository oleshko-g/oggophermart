// Package balance declares the user storage interface
package balance

type Storage interface {
	RetrieveUserBalance(userID string) (int, error)
	SaveUserTransaction(userID string, amount int) (error)
}
