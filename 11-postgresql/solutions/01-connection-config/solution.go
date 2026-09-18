package poolconfig_solution

import (
	"database/sql"
	"fmt"
	"time"
)

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ValidateAndApply validates pool limits and applies them to *sql.DB.
func ValidateAndApply(db *sql.DB, cfg PoolConfig) error {
	if cfg.MaxOpenConns <= 0 {
		return fmt.Errorf("max open conns must be > 0, got %d", cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns <= 0 || cfg.MaxIdleConns > cfg.MaxOpenConns {
		return fmt.Errorf("max idle conns must be > 0 and <= max open conns (%d), got %d", cfg.MaxOpenConns, cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime <= 0 {
		return fmt.Errorf("conn max lifetime must be > 0, got %v", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime <= 0 || cfg.ConnMaxIdleTime > cfg.ConnMaxLifetime {
		return fmt.Errorf("conn max idle time must be > 0 and <= conn max lifetime (%v), got %v", cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return nil
}
