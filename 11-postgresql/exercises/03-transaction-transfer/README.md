# Exercise 03: Atomic Transaction Transfer (Tier 3 - Hard)

## 🎯 Problem Statement
In financial and inventory systems, balance mutations across multiple entities must occur atomically inside an ACID transaction. If an error occurs halfway through, changes must be rolled back cleanly.

Implement `AtomicTransfer`:
```go
type TxExecutor interface {
    GetBalance(ctx context.Context, accountID string) (float64, error)
    UpdateBalance(ctx context.Context, accountID string, delta float64) error
    Commit() error
    Rollback() error
}

func AtomicTransfer(ctx context.Context, tx TxExecutor, fromID, toID string, amount float64) error
```

### Requirements & Business Rules
1. If `amount <= 0`, return `ErrInvalidAmount`.
2. If `fromID == toID`, return `ErrSameAccount`.
3. Ensure `defer tx.Rollback()` is registered immediately.
4. Retrieve sender's balance using `tx.GetBalance`.
5. If `balance < amount`, return `ErrInsufficientFunds`.
6. Debit sender: `tx.UpdateBalance(ctx, fromID, -amount)`.
7. Credit receiver: `tx.UpdateBalance(ctx, toID, amount)`.
8. Commit transaction: `tx.Commit()`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `AtomicTransfer`.
3. Run tests:
   ```powershell
   go test -v ./11-postgresql/exercises/03-transaction-transfer
   ```
