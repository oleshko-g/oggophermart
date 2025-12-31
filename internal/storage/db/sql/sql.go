// Package sql is the impementation of
package sql
import _ "github.com/oleshko-g/oggophermart/internal/storage"
import "github.com/oleshko-g/oggophermart/internal/storage/db"
import "github.com/oleshko-g/oggophermart/internal/storage/db/sql/schema"
import "database/sql"
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
