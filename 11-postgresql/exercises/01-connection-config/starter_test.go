package poolconfig

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"
)

type dummyDriver struct{}

func (d dummyDriver) Open(name string) (driver.Conn, error) {
	return nil, nil
}

func init() {
	sql.Register("dummy_pool_test", dummyDriver{})
}

func TestValidateAndApply(t *testing.T) {
	db, err := sql.Open("dummy_pool_test", "")
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

	err = ValidateAndApply(db, validCfg)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: ValidateAndApply is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("expected valid config to succeed, got %v", err)
	}

	// Test invalid cases
	invalidCases := []struct {
		name string
		cfg  PoolConfig
	}{
		{
			name: "zero max open conns",
			cfg: PoolConfig{
				MaxOpenConns:    0,
				MaxIdleConns:    0,
				ConnMaxLifetime: 10 * time.Minute,
				ConnMaxIdleTime: 2 * time.Minute,
			},
		},
		{
			name: "idle conns exceeds open conns",
			cfg: PoolConfig{
				MaxOpenConns:    10,
				MaxIdleConns:    20,
				ConnMaxLifetime: 10 * time.Minute,
				ConnMaxIdleTime: 2 * time.Minute,
			},
		},
		{
			name: "idle time exceeds max lifetime",
			cfg: PoolConfig{
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 10 * time.Minute,
			},
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateAndApply(db, tc.cfg); err == nil {
				t.Errorf("expected error for invalid config %s, got nil", tc.name)
			}
		})
	}
}
