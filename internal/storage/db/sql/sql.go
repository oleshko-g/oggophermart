// Package sql is the impementation of
package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // revive:disable-line:blank-imports registers the postgres driver
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql/query"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql/schema"
	"github.com/oleshko-g/oggophermart/internal/storage/errors"
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

// RetrieveUserBalance retrieves current user's balance and the amount withdrawn by their userID or an error
func (s *Storage) RetrieveUserBalance(userID string) (currentBalance, withdrawn int, err error) {
	return 0, 0, nil
}

// SaveUserTransaction saved the user's transaction by the following logic:
//   - a) If the amount is postive then it's an accrual
//   - b) if the amount is negative then it's a withdrawl
func (s *Storage) SaveUserTransaction(userID string, amount int) error { return nil }

// RetrieveUser retrieves a single user by their id
func (s *Storage) RetrieveUser(id string) error { return nil }

// StoreUser stores the user by their name and their hashed password.
//   - name MUST be unique
func (s *Storage) StoreUser(ctx context.Context, login, hashedPassword string) (err error) {
	result, err := s.db.ExecContext(ctx, query.SelectUserIDByLogin, login)
	if err != nil {
		return err
	}

	num, err := result.RowsAffected()
	if num != 0 {
		return fmt.Errorf("%w: user login", errors.ErrAlreadyExists)
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	result, err = s.db.ExecContext(ctx,
		query.InsertUser,
		newUserID,
		login,
		hashedPassword,
		time.Now().UTC(), // creatd_at
		time.Now().UTC(), // updated_at
	)
	if err != nil {
		return err
	}

	num, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if num != 1 {
		return fmt.Errorf("%w: expected to affect 1 row, affected %d", errors.ErrNoAffect, num)
	}

	return nil
}

var _ storage.User = (*Storage)(nil)
var _ storage.Balance = (*Storage)(nil)
