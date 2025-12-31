// Package sql is the impementation of
package sql

import (
	"database/sql"

	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql/schema"
)

// New configures and open a new connection to the db and returns a [Storage] or an error
func New(c *db.Config) (s *Storage, err error) {
	database, err := sql.Open(c.DSN().DriverName.String(), c.DSN().String())
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	if err = schema.Up(c.DSN().DriverName, database); err != nil {
		return
	}

	return &Storage{
		db: database,
	}, nil
}

// Storage represents an internal implementation of [sql.DB]
type Storage struct {
	db *sql.DB
}

func (s *Storage) RetrieveUserBalance(userID string) (int, error)               { return 0, nil }
func (s *Storage) SaveUserTransaction(userID string, amount int) error          { return nil }
func (s *Storage) RetrieveUser(id string) error                                 { return nil }
func (s *Storage) StoreUser(name, hashedPassword string) (id string, err error) { return "", nil }

var _ storage.User = (*Storage)(nil)
var _ storage.Balance = (*Storage)(nil)
