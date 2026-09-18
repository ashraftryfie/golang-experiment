# Stage 11: PostgreSQL & Database Persistence

Welcome to **Stage 11** of your Go engineering journey. In this stage, you master database programming with Go's standard `database/sql` package, production connection pool tuning, context-aware queries, parameterized query safety, ACID transactions, and database migration architecture.

---

## 🧭 Mental Model

### 1. `*sql.DB` Is a Connection Pool, Not a Connection
In Go, `sql.Open("postgres", dsn)` does not open an actual network socket immediately; it initializes a lazy thread-safe connection pool (`*sql.DB`):
```go
db, err := sql.Open("pgx", dsn)
if err != nil {
    return err
}
// Validate connection with context timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
if err := db.PingContext(ctx); err != nil {
    return err
}
```

### 2. Connection Pool Sizing Best Practices
Default Go settings keep unlimited open connections and 2 idle connections, which causes rapid connection churning and exhaustion under load:
```go
db.SetMaxOpenConns(25)                 // Max active + idle connections
db.SetMaxIdleConns(25)                 // Keep pool warm
db.SetConnMaxLifetime(15 * time.Minute) // Retire old connections gracefully
db.SetConnMaxIdleTime(5 * time.Minute)  // Reclaim unused idle sockets
```

### 3. Context-Aware Execution & Leak Prevention
Never use bare `db.Query()` or `db.Exec()`. Always pass a `context.Context` to allow cancellation on client disconnects and timeouts:
```go
rows, err := db.QueryContext(ctx, "SELECT id, title FROM tasks WHERE status = $1", status)
if err != nil {
    return nil, err
}
defer rows.Close() // MANDATORY: Releases connection back to pool!

for rows.Next() {
    var id, title string
    if err := rows.Scan(&id, &title); err != nil {
        return nil, err
    }
}
// MANDATORY: Check for errors encountered during row iteration!
if err := rows.Err(); err != nil {
    return nil, err
}
```

### 4. Canonical Transaction Pattern
In Go, calling `defer tx.Rollback()` is standard and safe: if `tx.Commit()` succeeds, the deferred rollback returns `sql.ErrTxDone` and is a harmless no-op:
```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback() // Safe no-op if committed

if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID); err != nil {
    return fmt.Errorf("debit failed: %w", err)
}

if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID); err != nil {
    return fmt.Errorf("credit failed: %w", err)
}

return tx.Commit() // Persists changes atomically
```

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Configure production connection pools with appropriate open/idle bounds.
- [x] Execute context-aware parameterized queries with SQL injection immunity.
- [x] Prevent connection pool starvation using `defer rows.Close()` and `rows.Err()`.
- [x] Execute multi-statement ACID transactions with bulletproof deferred rollbacks.
- [x] Structure database migrations tracked by schema versioning tables.

---

## 🗂️ Stage Structure

```
11-postgresql/
├── README.md                           # Curriculum & architectural patterns
├── examples/
│   ├── connection_pool/main.go         # Connection pool tuning & health check
│   └── transaction_rollback/main.go    # Safe ACID transfer with deferred rollback
├── exercises/
│   ├── 01-connection-config/           # Connection pool builder & validator
│   ├── 02-query-scanner/               # Context-aware row scanner & resource cleanup
│   └── 03-transaction-transfer/        # Atomic bank transfer with balance check
├── solutions/
│   ├── 01-connection-config/           # Reference implementation
│   ├── 02-query-scanner/               # Reference implementation
│   └── 03-transaction-transfer/        # Reference implementation
└── challenges/
    └── migration-runner/               # Transactional schema migration engine
```

---

## 🧪 Verification Commands

```powershell
# Run all tests in this stage
go test -v ./11-postgresql/...

# Vet and check format
go vet ./11-postgresql/...
go fmt ./11-postgresql/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 10: REST API Engineering](../10-rest-api/README.md).
- **Next Stage**: [Stage 12: Testing & Benchmarking](../12-testing/README.md) (Table tests, `httptest`, race detector, benchmarks, fuzzing).
