
// Package user declares the user storage interface
package user

type Storage interface {
	RetrieveUser(id string) error
	StoreUser(name, hashedPassword string) (id string, err error)
}
