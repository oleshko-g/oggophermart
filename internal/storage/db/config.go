// Package db is the base package for db servers
package db

import (
	"net/url"

	storageErrors "github.com/oleshko-g/oggophermart/internal/storage/errors"
)

// Config represents a config of an SQL database
type Config struct {
	dataSource
}

// DSN returns a pointer to the [flag.Value] to set the database source name
func (c *Config) DSN() *dataSource { // revive:disable-line:unexported-return provides the interface to the caller
	return &c.dataSource
}

// dataSource represent a valid Data Source
type dataSource struct {
	name string
	DriverName
	Source string
}

// Set parses s and sets [DSN] and [Driver] or returns an error
func (d *dataSource) Set(s string) error {
	switch d.Source {
	case "":
		d.Source = "env var" // assumes that env var is set upon the first call to [Set]
	case "env var":
		d.Source = "command line flag" // assumes that flags are always parsed last
	}

	url, err := url.Parse(s)
	if err != nil {
		return err
	}

	if url.Scheme != string(DriverNamePostgres) && url.Scheme != string(DriverNamePostgreSQL) {
		return storageErrors.ErrUnsupportedDataSource
	}

	// there's only "postgres" SQL driver
	d.DriverName = DriverName(DriverNamePostgres)

	d.name = url.String()

	return nil
}

func (d *dataSource) String() string {
	return d.name
}

// DriverName is a valid database driver name
type DriverName string

func (d DriverName) String() string {
	return string(d)
}

// Supported database drivers
const (
	DriverNamePostgres   DriverName = "postgres"
	DriverNamePostgreSQL DriverName = "postgresql"
)
