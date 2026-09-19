# Exercise 03: Service Layer & Dependency Injection

## Objective
Implement a robust, testable **Service Layer** (Use-Case orchestrator) using constructor-based **Dependency Injection** in Go, coordinating repository transactions, invariant checks, compensation on error, and event notification.

## Concepts Covered
- **Service Layer**: Coordinates business workflows, transactions, domain entity mutations, and outbound notifications.
- **Constructor Injection**: In Go, dependency injection is performed via explicit constructor functions (e.g. `NewTransferService(repo, notifier)`), avoiding magical reflection or runtime container frameworks.
- **Error Handling & Compensation**: If a subsequent persistence step fails during multi-entity coordination, previous changes must be rolled back (compensating transaction) or wrapped with meaningful context (`fmt.Errorf("...: %w", err)`).

## Requirements
1. **Errors**:
   - `ErrInvalidAmount`: Amount <= 0.
   - `ErrSameAccount`: Sender and receiver account IDs are identical.
   - `ErrInsufficientFunds`: Sender balance is less than transfer amount.
   - `ErrAccountNotFound`: Account cannot be found in repository.

2. **Domain (`Account`)**:
   - `NewAccount(id string, balance int64) *Account`
   - `Debit(amount int64) error`
   - `Credit(amount int64)`

3. **Ports**:
   - `AccountRepository`: `Get(ctx context.Context, id string) (*Account, error)` and `Update(ctx context.Context, a *Account) error`
   - `Notifier`: `NotifyTransfer(ctx context.Context, toID string, amount int64) error`

4. **Service (`TransferService`)**:
   - `NewTransferService(repo AccountRepository, notifier Notifier) *TransferService`
   - `ExecuteTransfer(ctx context.Context, fromID, toID string, amount int64) error`:
     - Validates input (`amount > 0`, `fromID != toID`).
     - Loads sender and receiver.
     - Debits sender and credits receiver.
     - Persists sender update. If fails, return error.
     - Persists receiver update. If fails, rollback sender (re-credit) and return error.
     - Dispatches notification to receiver.
