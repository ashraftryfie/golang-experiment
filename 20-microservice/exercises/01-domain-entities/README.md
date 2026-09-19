# Exercise 01: Domain Entities & Value Objects

## Objective
Implement pure Domain-Driven Design (DDD) domain models in Go, separating immutable **Value Objects** from mutable **Entities**, protecting domain invariants, and returning typed domain errors.

## Concepts Covered
- Value Objects: Immutable structures identified solely by their data (e.g. `Money`).
- Entities: Domain concepts identified by unique identities and state transitions (e.g. `Wallet`).
- Domain Invariants: Business rules enforced inside constructors and mutating methods.
- Typed Domain Errors: Domain-specific sentinel errors (`ErrNegativeAmount`, `ErrCurrencyMismatch`, `ErrInsufficientBalance`).

## Requirements
1. **`Money` (Value Object)**:
   - Fields should be unexported (`cents int64`, `currency string`).
   - `NewMoney(cents int64, currency string) (Money, error)`: Validates `cents >= 0` and `currency != ""`.
   - `Add(other Money) (Money, error)`: Adds two amounts; returns `ErrCurrencyMismatch` if currencies differ.
   - `Subtract(other Money) (Money, error)`: Subtracts other amount; returns `ErrCurrencyMismatch` if currencies differ, or `ErrInsufficientBalance` if resulting cents < 0.
   - `Equals(other Money) bool`: Returns true if both cents and currency match.
   - `String() string`: Formats as `"10.50 USD"` (formatted decimal with two decimal places and currency).

2. **`Wallet` (Entity)**:
   - Encapsulated fields (`id`, `ownerID`, `balance`).
   - `NewWallet(id, ownerID string, initialBalance Money) (*Wallet, error)`: Rejects empty IDs with `ErrInvalidID`.
   - `Deposit(m Money) error`: Enforces currency match and positive deposit amount.
   - `Withdraw(m Money) error`: Enforces currency match and balance check; updates balance.
