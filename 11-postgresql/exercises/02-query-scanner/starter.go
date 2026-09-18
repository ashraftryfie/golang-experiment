package scanner

import (
	"errors"
)

var ErrNotImplemented = errors.New("TODO: implement ScanUsers")

type UserRecord struct {
	ID       int64
	Email    string
	IsActive bool
}

// RowIterator models the lifecycle of *sql.Rows.
type RowIterator interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// ScanUsers scans all rows into UserRecord slices, guaranteeing Close and Err checks.
func ScanUsers(rows RowIterator) ([]UserRecord, error) {
	// TODO: defer rows.Close(), loop Next(), Scan(), check rows.Err()
	return nil, ErrNotImplemented
}
