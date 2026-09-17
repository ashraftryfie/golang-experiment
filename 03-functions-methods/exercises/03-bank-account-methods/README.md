# Exercise 03: Bank Account Method Set (Tier 3 - Hard)

## 🎯 Problem Statement
Design an immutable-ledger based `BankAccount` custom type demonstrating Go's value and pointer receiver rules, error checking, and method sets.

### Requirements
1. `NewAccount(owner string, initialDeposit float64) (*BankAccount, error)`:
   - Rejects empty owner names (`ErrInvalidOwner`).
   - Rejects initial deposits `< 0` (`ErrNegativeAmount`).
2. `(a *BankAccount) Deposit(amount float64) error`:
   - Validates `amount > 0` (`ErrNegativeAmount`).
   - Appends a transaction record and updates current balance.
3. `(a *BankAccount) Withdraw(amount float64) error`:
   - Validates `amount > 0` (`ErrNegativeAmount`).
   - Validates that balance has sufficient funds (`ErrInsufficientFunds`).
   - Appends a withdrawal transaction and decrements balance.
4. `(a BankAccount) Balance() float64`:
   - Value receiver: safely reads current balance without risk of state mutation.
5. `(a BankAccount) Statement() []Transaction`:
   - Returns an isolated copy of all transactions (protecting internal ledger from caller mutation).

```go
type TransactionType string

const (
    TxDeposit  TransactionType = "DEPOSIT"
    TxWithdraw TransactionType = "WITHDRAW"
)

type Transaction struct {
    Type   TransactionType
    Amount float64
}
```

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement methods ensuring transactions cannot mutate internal state.
3. Run tests:
   ```powershell
   go test -v ./03-functions-methods/exercises/03-bank-account-methods
   ```
