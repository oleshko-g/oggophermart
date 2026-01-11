// Package sql is the impementation of
package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // revive:disable-line:blank-imports registers the postgres driver
	genDBSQL "github.com/oleshko-g/oggophermart/internal/gen/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql/schema"
	storageErrors "github.com/oleshko-g/oggophermart/internal/storage/errors"
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

	queries := genDBSQL.New(database)

	return &Storage{
		db:      database,
		queries: queries,
	}, nil
}

// Storage represents an internal implementation of [sql.DB]
type Storage struct {
	db      *sql.DB
	queries *genDBSQL.Queries
}

var _ storage.User = (*Storage)(nil)
var _ storage.Balance = (*Storage)(nil)

// RetrieveUserBalance retrieves current user's balance and the amount withdrawn by their userID or an error
func (s *Storage) RetrieveUserBalance(ctx context.Context, userID uuid.UUID) (currentBalance, withdrawn int, err error) {
	return 0, 0, nil
}

// SaveUserTransaction saved the user's transaction by the following logic:
//   - a) If the amount is positive then it's an accrual
//   - b) if the amount is negative then it's a withdrawl
func (s *Storage) SaveUserTransaction(ctx context.Context, userID uuid.UUID, amount int) error {
	return nil
}

// RetrieveUser retrieves a user id by their login
func (s *Storage) RetrieveUser(ctx context.Context, login string) (userID uuid.UUID, err error) {
	userID, err = s.queries.SelectUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, storageErrors.ErrNotFound
		}
		return uuid.UUID{}, err
	}
	return userID, nil
}

// StoreUser stores the user by their name and their hashed password.
//   - name MUST be unique
func (s *Storage) StoreUser(ctx context.Context, login, hashedPassword string) (err error) {
	_, err = s.RetrieveUser(ctx, login)
	if err != nil {
		if errors.Is(err, storageErrors.ErrNotFound) {
			goto newUser
		}
		return err
	} else {
		return storageErrors.ErrAlreadyExists
	}

newUser:
	newUserID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	result, err := s.queries.InsertUser(ctx,
		genDBSQL.InsertUserParams{
			ID:             newUserID,
			Login:          login,
			HashedPassword: hashedPassword,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC()},
	)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num != 1 {
		return fmt.Errorf("%w: expected to affect 1 row, affected %d", storageErrors.ErrNoAffect, num)
	}

	return nil
}

func (s *Storage) StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, orderStatus string, createdAt time.Time) error {

	newOrderID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	res, err := s.queries.InsertOrder(ctx,
		genDBSQL.InsertOrderParams{
			ID:        newOrderID,
			UserID:    userID,
			Number:    orderNumber,
			Status:    orderStatus,
			CreatedAt: createdAt,
		})
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return storageErrors.ErrAlreadyExists
	}

	return nil
}
func (s *Storage) RetreiveOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error {
	// s.queries.Se
	return nil
}

func (s *Storage) RetreiveUserPassword(ctx context.Context, login string) (hashedPassword string, err error) {
	hashedPassword, err = s.queries.SelectUserHashedPasswordByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storageErrors.ErrNotFound
		}
		return "", err
	}
	return hashedPassword, nil
}

func (s *Storage) RetreiveOrderUser(ctx context.Context, orderNumber string) (userID uuid.UUID, err error) {
	userID, err = s.queries.SelectUserIDByOrderNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, storageErrors.ErrNotFound
		}

		return uuid.UUID{}, err
	}
	return userID, nil
}

func (s *Storage) RetrieaveUserOrders(ctx context.Context, userID uuid.UUID) (userOrders []genDBSQL.SelectOrdersByUserIDRow, err error) {
	rows, err := s.queries.SelectOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func (s *Storage) Retrieve(ctx context.Context, userID uuid.UUID) (genDBSQL.UserBalance, error) {
	userBalance, err := s.queries.SelectBalanceByUserID(ctx, userID)
	if err != nil {
		return genDBSQL.UserBalance{}, err
	}
	return userBalance, nil
}
