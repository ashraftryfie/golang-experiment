package poolconfig

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("TODO: implement ValidateAndApply")

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ValidateAndApply validates pool limits and applies them to *sql.DB.
func ValidateAndApply(db *sql.DB, cfg PoolConfig) error {
	// TODO: Validate cfg constraints and call db.SetMaxOpenConns, db.SetMaxIdleConns, etc.
	return ErrNotImplemented
}
