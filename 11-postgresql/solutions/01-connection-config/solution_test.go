package poolconfig_solution

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
)

type dummyDriver struct{}

func (d dummyDriver) Open(name string) (driver.Conn, error) {
	return nil, nil
}

func init() {
	sql.Register("dummy_pool_sol", dummyDriver{})
}

func TestValidateAndApplySolution(t *testing.T) {
	db, err := sql.Open("dummy_pool_sol", "")
	if err != nil {
		t.Fatalf("failed to open dummy db: %v", err)
	}
	defer db.Close()

	validCfg := PoolConfig{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	if err := ValidateAndApply(db, validCfg); err != nil {
		t.Fatalf("expected valid config to succeed, got %v", err)
	}

	invalidCfg := PoolConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    20, // exceeds max open
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	if err := ValidateAndApply(db, invalidCfg); err == nil {
		t.Errorf("expected error for invalid config, got nil")
	}
}
