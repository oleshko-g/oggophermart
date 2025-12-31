package storage

// User declares the storage interface for the user service
type User interface {
	RetrieveUser(id string) error
	StoreUser(name, hashedPassword string) (id string, err error)
}
