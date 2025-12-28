// Package errors is the internal storage package
package errors //revive:disable:var-naming

import (
	"errors"
)

var (
	// ErrUnsupportedDataSource is the error returned when data source is not supported by the storage implementation
	ErrUnsupportedDataSource = errors.New("unsupported data source")
)
