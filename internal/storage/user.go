package storage

// Storage user declares the user storage interface
type User interface {
	RetrieveUser(id string) error
	StoreUser(name, hashedPassword string) (id string, err error)
}
