package scanner_solution

import (
	"fmt"
)

type UserRecord struct {
	ID       int64
	Email    string
	IsActive bool
}

type RowIterator interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// ScanUsers scans all rows into UserRecord slices, guaranteeing Close and Err checks.
func ScanUsers(rows RowIterator) ([]UserRecord, error) {
	defer rows.Close()

	users := make([]UserRecord, 0)

	for rows.Next() {
		var u UserRecord
		if err := rows.Scan(&u.ID, &u.Email, &u.IsActive); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}
