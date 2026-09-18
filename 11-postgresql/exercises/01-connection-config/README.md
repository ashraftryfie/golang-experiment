# Exercise 01: Connection Pool Configuration & Validation (Tier 1 - Easy)

## 🎯 Problem Statement
Misconfigured database connection pools cause connection storms, socket exhaustion, or high latency.

Implement `ValidateAndApply`:
```go
type PoolConfig struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func ValidateAndApply(db *sql.DB, cfg PoolConfig) error
```

### Validation Rules
1. `MaxOpenConns` must be $> 0$.
2. `MaxIdleConns` must be $> 0$ and $\le$ `MaxOpenConns`.
3. `ConnMaxLifetime` must be $> 0$.
4. `ConnMaxIdleTime` must be $> 0$ and $\le$ `ConnMaxLifetime`.
5. If invalid, return a descriptive error.
6. If valid, apply settings to `db` using `db.SetMaxOpenConns`, `db.SetMaxIdleConns`, `db.SetConnMaxLifetime`, `db.SetConnMaxIdleTime`, and return `nil`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `ValidateAndApply`.
3. Run tests:
   ```powershell
   go test -v ./11-postgresql/exercises/01-connection-config
   ```
