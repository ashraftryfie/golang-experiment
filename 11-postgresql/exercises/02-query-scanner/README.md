# Exercise 02: Row Scanner & Resource Lifecycle (Tier 2 - Medium)

## 🎯 Problem Statement
Failing to close database rows leads to socket leaks and connection pool starvation. Failing to check `rows.Err()` silently ignores data transmission truncations.

Implement `ScanUsers`:
```go
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

func ScanUsers(rows RowIterator) ([]UserRecord, error)
```

### Requirements
1. Always call `defer rows.Close()` at function entry.
2. Iterate using `rows.Next()`.
3. Scan columns into `UserRecord` fields: `&u.ID`, `&u.Email`, `&u.IsActive`.
4. If `Scan` returns an error, exit immediately and return the error.
5. After the iteration loop, inspect `rows.Err()`. If non-nil, return it.
6. Return a non-nil slice `[]UserRecord{}` when no rows are found.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `ScanUsers`.
3. Run tests:
   ```powershell
   go test -v ./11-postgresql/exercises/02-query-scanner
   ```
