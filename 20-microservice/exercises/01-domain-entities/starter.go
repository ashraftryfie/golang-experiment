package domainentities

import (
	"errors"
)

var (
	ErrNegativeAmount      = errors.New("negative amount not allowed")
	ErrInvalidCurrency     = errors.New("currency code cannot be empty")
	ErrCurrencyMismatch    = errors.New("currency mismatch")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidID           = errors.New("id cannot be empty")
	ErrNotImplemented      = errors.New("not implemented yet")
)

// Money is an immutable Value Object
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	// TODO: Implement validation and return Money
	return Money{}, ErrNotImplemented
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	// TODO: Verify currencies match, return sum
	return Money{}, ErrNotImplemented
}

func (m Money) Subtract(other Money) (Money, error) {
	// TODO: Verify currencies match and result >= 0, return difference
	return Money{}, ErrNotImplemented
}

func (m Money) Equals(other Money) bool {
	// TODO: Return true if cents and currency match
	return false
}

func (m Money) String() string {
	// TODO: Format as "10.50 USD"
	return ""
}

// Wallet is an Aggregate Entity
type Wallet struct {
	id      string
	ownerID string
	balance Money
}

func NewWallet(id, ownerID string, initialBalance Money) (*Wallet, error) {
	// TODO: Validate ID and ownerID, instantiate Wallet
	return nil, ErrNotImplemented
}

func (w *Wallet) ID() string {
	return w.id
}

func (w *Wallet) OwnerID() string {
	return w.ownerID
}

func (w *Wallet) Balance() Money {
	return w.balance
}

func (w *Wallet) Deposit(m Money) error {
	// TODO: Enforce invariants and update balance
	return ErrNotImplemented
}

func (w *Wallet) Withdraw(m Money) error {
	// TODO: Enforce invariants, check balance, and update balance
	return ErrNotImplemented
}
